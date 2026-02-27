This document is a 'stream of conciousness' of my thought process as I built this tool. It is a little verbose as I tried as much as possible to capture the 'why', not just the 'how'.

## AI Usage
I did not use any AI for the creation of this tool, whilst I am not averse to the usage of AI, if the intention of the test is to identify my coding style and knowledge, the usage of AI feels like it would obfuscate this to an extent.

## Language choice
I use Go and Rust in my current day job, but more Go than rust. Whilst I switch frequently at the moment I tend to 'think more in Go' and can churn out code quicker in it, so it was the logical choice for the language on this occasion.

## Overall code design
The design was to carry out a simple byte by byte shift. I would read the file in and write it out in chunks, bit shifting each chunk as I go. For left shift I would go 'left to right' through the file. In both cases I have to remember the byte that will 'fall off the end' and make sure it is written to the start/end of the file.

## ShiftByte
The first function created was ShiftByte the intention here is to create a simple byte shift function that can be re-used in a tight loop for each chunk.

I took a purely TDD approach, creating first just the LEFT/RIGHT 'enum' and the layout of the function, returning 0,0. Then creating the table driven test and the first simple functions. I then proceeded to come up with some inital test cases, adding code as I went to ensure it passed. Later on I added the overflow output and a structure for the result to improve readability.

After exhausting the manual test cases I went for the fuzz function just in case there was edge cases I had missed.

I could add some belt and braces input validation here, in particular ensuring that the 'overflow' byte is either 00000001 or 10000000, but as this is an internal only function it felt overkill.

## ShiftChunk
Shift chunk really looks like a scaled up version of shift byte, it too takes an input overflow byte (from the previous chunk if one exists) and it will call shift on each input byte as it goes along, writing to an output chunk and passing the previous overflow to the next call.

Testing followed a very similar vein to ShiftByte.

I'm not entirely happy with the duplicated code in the ShiftChunk method, I possibly could have created an interator that can go left or right on a chunk, then the duplicated code could be unified.

## ShiftFile
The shift file function is a bit of a convoluted beast, as it has to handle shifting the file left and right.
If we are shifting left then we first determine the overflow bit that would be shifted off the left from the start of the file. Then we seek the end of the file and enter a loop where we shift left by 'chunk' bytes, do a 'chunk' shift, write to the output file, rinse and repeat.
Shifting right is simpler, we first determine what overflow bit would shift off the end of the file, then we go back to the start. After that our tight loop just reads a chunk, shifts it, writes it back to disk, and so on.

I pulled some code out into local functions as I went along just to keep the ShiftFile function getting to overwhelming. I suspect I could have broken up those functions a bit further, but it felt like I was hitting diminishing returns.

Note that the chunk size was hard coded, the size could be larger and in a real world I'd likely make the chunk size configurable as it can then be tuned based on available memory, etc. But ultimately that isn't particularly important to the proof of the algorithm itself working.

Again TDD was followed. I had taken close to 2 hours at this point, so I haven't added a fuzz function. Once the main was created (below) I did do some round trip test (left/right and right/left) on a few files to confirm it all worked.

## Main.go
Finally the main file was created, in a real CLI tool I would likely use cobra+viper library for handling command line options but I wanted to keep it simple here, so it's all just manual parsing of the command line.

## Notes
All code was stored in a local git repo, with regular commits (these commits were lost when I copied over to github).

### Enhancements
I feel this is enough to show the basic operation of the algorithm, but there is some enhancements that may have been possible if there were performance requirements specified:

* Currently for each call to ShiftChunk I create a new output chunk, this leads a to clean API but will result in a lot of allocations, which will hit performance in the tight loop of ShiftFile. Alternatives are to pass in an output buffer, this would allow ShiftFile to allocate the buffer once and re-use it. Alternatively we could write back to the input buffer, thus reducing even the single output buffer allocation.
* Rather than truncating the input file I could have initialised it to the correct length, I believe that would mean that an in place shift (file_in == file_out) would work (since file_out would already be the correct length). In fact thinking about it, I could just check if the strings were equal, and then open file_in as r/w and pass it to ShiftFile function twice.
* Mmap the files, I could just use a single mmaped chunk in this case.
* Parallelise the operation. This would require co-ordination of the overflow bytes, but should be doable as the last/first overflow byte could be written after the chunk has been completed. Would of course use more memory, although this would work will with mmap (as I could pass slices to the original mmaped file to each parallel go routine).
* Use uint64 rather than uint8, this should be quicker but would need work to handle endianness.

### Tidy ups
The overflow passed in to ShiftChunk is a byte, it might be neater for the API for it to be a bool indicating the overflow bits presence or otherwise (similarly in the output). I duplicate the overflow bit calculation in ShiftFile to calculate the 'bit that falls out the end' I'm not overly happy with that and that could potentially be rolled into the shift package as a standalone function on a single byte. Having said that again for the purposes of showing off the basic algorithm and code style I think it's adequate, in a production environment I would likely be asking questions about how the tool and the underlying library was to be used and tweak the API appropriately.
Finally I truncate the output file, this means if file_in == file_out we'll fail and destroy the original file. This is fixable by determining if they are the same file and opening file_in as r/w and passing it through to ShiftFile as both the input and the output.