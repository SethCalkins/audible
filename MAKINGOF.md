## Making of: Audible header parser

### A note upfront

The reason I did all this, is because I originally wanted to find out if I
could work around the DRM stuff in Audible audio books. So far this has not
been successful, but at least we got a header parser library out of it.

While writing this, I found there is a python library which interprets .aa file
header data. While I could have used it as a reference to make my life easier,
I decided to ignore it and try to figure it all out by hand. Mostly because
it's a lot more fun and educational this way.


### Setting the scene

So here I am, watching one of my usual goto [youtube channels][vt] and
towards the end I am presented with yet another plug for something called
[Audible.com][ad].

[vt]: https://www.youtube.com/user/1veritasium
[ad]: http://www.audible.com/

I've seen people do these plugs a billion times by now and never really
cared (I am rather allergic to anything advertising related, no matter how
benign). But for some reason, this time, I figured I give in and check it out.
What could possibly go wrong?

I reluctantly sign up for an account and am immediately presented with a free,
introductory audio book. "Well, isn't that nice?", I thought. "Let's download
that stuff and give it a whirl!". A few minutes later I was the proud owner
of a shiny 1.2MB `.aa` file, sitting on my disk. I never heard of this format,
so I was hoping [VLC][vlc] or [mplayer][mpl] might have. Otherwise I would not
be able to listen to whatever goodies where awaiting me inside.

[vlc]: http://www.videolan.org/vlc/index.html
[mpl]: https://mplayerhq.hu

Unfortunatly, neither of these had any idea what to do with the file, so
I had to go for a different approach. I hopped onto Google and punched
in a few queries to see what came up. Almost immediately I realized this was
going to be more difficult than I had hoped. The first results covered
questions from others, asking how to convert `.aa` to `.mp3` and how to
circumvent the DRM parts. This did not inspire me with confidence. The gist of
the 'solution' was simple:

* Open the .aa file in some Audible-approved software like iTunes.
* Burn/export the audio as an audio CD image.
* Use a CD ripper to convert to mp3.

Needless to say this won't do at all. Not in the least because I don't do
iTunes. But also because, frankly, it's stupid that one needs to jump through
hoops like this, just to be able to listen to something you should be fully
entitled to listen to in any way you like.

Considering I did not really have anything better to do and always wanted to
dive deeper into the 'reverse engineering' business, I figured this was my
queue to break out the coffee. Instead of listening to a great adventure in
some book, I would create my own. No doubt riddled with dragons and dangerous
traps along the way, but this did not phase me in the least. "I am ready!",
I thought; "Bring it on!".


### Testing the water

Let me start by making it clear that I have no idea what I am doing here.
I am simply poking the hornet's nest with a stick to see what happens and I'll
take it from there. While I have 14 years of programming experience under my
belt (it's a big belt) and thus have encountered a few file format
specifications, I am very much unaccustomed to divining said specification out
of thin air.

The first thing I do with any unknown binary file, is cast upon it some magic
missles in the shape of [binwalk](https://github.com/devttys0/binwalk). This
tool is intended to examine firmware images, but I found it can yield at
least some useful information on pretty much any binary file.

```
$ binwalk -B data/file.aa
DECIMAL       HEXADECIMAL     DESCRIPTION
--------------------------------------------------------------------------------
1212700       0x12811C        JPEG image data, JFIF standard  1.02
```

Well look at that. It's a JPEG image. Let's extract that and see what it
contains:

	$ dd if=data/file.aa of=data/file.jpeg bs=1 skip=1212700

Opening the resulting jpeg presented me with a quaint little logo, showing
the title of the book. Not particularly helpful, but at least now we have
2.6kB of data less to worry about. The rest of the file is going to take a
bit more effort it seems.

There has to be some kind of file header which is not encrypted. Perhaps even
some plain text metadata. Out of curiosity, let's just extract the first 4kB
of data from the file and examine its contents. We have no way of knowning if
this covers the entire header, only part of it or if there's a header at all.
But we have to start somewhere and this is as good a place as any to begin
digging for treasure.

```
$ dd if=data/file.aa of=data/header.bin bs=1024 count=4
$ hexdump -C data/header.bin
00000000  00 12 8b 6e 57 90 75 36  00 00 00 0c 00 00 00 00  |...nW.u6........|
00000010  00 00 00 00 00 00 00 00  00 12 8b 6e 00 00 00 01  |...........n....|
00000020  00 00 00 a0 00 00 00 18  00 00 00 02 00 00 00 b8  |................|
00000030  00 00 03 98 00 00 00 03  00 00 04 50 00 00 00 4a  |...........P...J|
00000040  00 00 00 04 00 00 04 9a  00 00 00 04 00 00 00 05  |................|
00000050  00 00 04 9e 00 00 00 08  00 00 00 06 00 00 04 a6  |................|
00000060  00 00 09 9e 00 00 00 07  00 00 0e 44 00 00 01 44  |...........D...D|
00000070  00 00 00 08 00 00 0f 88  00 00 00 04 00 00 00 09  |................|
00000080  00 00 0f 8c 00 00 00 00  00 00 00 0a 00 00 0f 8c  |................|
00000090  00 12 71 88 00 00 00 0b  00 12 81 14 00 00 0a 5a  |..q............Z|
000000a0  00 01 00 00 57 90 75 36  00 00 00 00 00 00 00 00  |....W.u6........|
000000b0  00 00 00 00 00 00 00 00  00 00 00 17 00 00 00 00  |................|
000000c0  0a 00 00 00 14 70 72 6f  64 75 63 74 5f 69 64 46  |.....product_idF|
000000d0  52 5f 50 52 4d 4f 5f 30  30 30 30 30 31 5f 6d 70  |R_PRMO_000001_mp|
000000e0  33 33 32 00 00 00 00 05  00 00 00 11 74 69 74 6c  |332.........titl|
000000f0  65 59 6f 75 72 20 46 69  72 73 74 20 4c 69 73 74  |eYour First List|
00000100  65 6e 00 00 00 00 08 00  00 00 0d 70 72 6f 76 69  |en.........provi|
00000110  64 65 72 41 75 64 69 62  6c 65 2c 20 49 6e 63 2e  |derAudible, Inc.|
00000120  00 00 00 00 07 00 00 00  0b 70 75 62 64 61 74 65  |.........pubdate|
00000130  30 31 2d 4a 41 4e 2d 32  30 30 30 00 00 00 00 0b  |01-JAN-2000.....|
00000140  00 00 00 2c 64 65 73 63  72 69 70 74 69 6f 6e 59  |...,descriptionY|
00000150  6f 75 72 20 41 75 64 69  62 6c 65 20 61 64 76 65  |our Audible adve|
00000160  6e 74 75 72 65 73 20 62  65 67 69 6e 20 72 69 67  |ntures begin rig|
00000170  68 74 20 68 65 72 65 2e  2e 2e 2e 00 00 00 00 10  |ht here.........|
00000180  00 00 00 59 6c 6f 6e 67  5f 64 65 73 63 72 69 70  |...Ylong_descrip|
00000190  74 69 6f 6e 59 6f 75 72  20 41 75 64 69 62 6c 65  |tionYour Audible|
000001a0  20 61 64 76 65 6e 74 75  72 65 73 20 62 65 67 69  | adventures begi|
000001b0  6e 20 72 69 67 68 74 20  68 65 72 65 2c 20 77 69  |n right here, wi|
000001c0  74 68 20 22 44 65 61 72  20 41 6d 61 6e 64 61 22  |th "Dear Amanda"|
000001d0  20 62 79 20 53 74 65 76  65 20 4d 61 72 74 69 6e  | by Steve Martin|
000001e0  2e 20 4c 69 73 74 65 6e  20 6e 6f 77 2e 00 00 00  |. Listen now....|
000001f0  00 09 00 00 00 0e 63 6f  70 79 72 69 67 68 74 28  |......copyright(|
00000200  43 29 41 75 64 69 62 6c  65 2e 63 6f 6d 00 00 00  |C)Audible.com...|
00000210  00 0b 00 00 00 11 73 68  6f 72 74 5f 74 69 74 6c  |......short_titl|
00000220  65 59 6f 75 72 20 46 69  72 73 74 20 4c 69 73 74  |eYour First List|
00000230  65 6e 00 00 00 00 0e 00  00 00 02 69 73 5f 61 67  |en.........is_ag|
00000240  67 72 65 67 61 74 69 6f  6e 6e 6f 00 00 00 00 08  |gregationno.....|
00000250  00 00 00 0e 74 69 74 6c  65 5f 69 64 46 52 5f 50  |....title_idFR_P|
00000260  52 4d 4f 5f 30 30 30 30  30 31 00 00 00 00 05 00  |RMO_000001......|
00000270  00 00 05 63 6f 64 65 63  6d 70 33 33 32 00 00 00  |...codecmp332...|
00000280  00 0a 00 00 00 0a 48 65  61 64 65 72 53 65 65 64  |......HeaderSeed|
00000290  31 35 32 38 32 38 35 34  39 35 00 00 00 00 0f 00  |1528285495......|
000002a0  00 00 05 45 6e 63 72 79  70 74 65 64 42 6c 6f 63  |...EncryptedBloc|
000002b0  6b 73 35 32 31 38 39 00  00 00 00 09 00 00 00 2a  |ks52189........*|
000002c0  48 65 61 64 65 72 4b 65  79 33 34 39 38 39 38 30  |HeaderKey3498980|
000002d0  31 38 36 20 31 39 38 38  37 30 31 39 30 20 33 30  |186 198870190 30|
000002e0  30 31 34 32 39 32 37 31  20 32 31 31 35 31 33 37  |01429271 2115137|
000002f0  30 38 38 00 00 00 00 0c  00 00 00 09 6c 69 63 65  |088.........lice|
00000300  6e 73 65 5f 6c 69 73 74  31 30 33 31 34 33 34 31  |nse_list10314341|
00000310  37 00 00 00 00 07 00 00  00 01 43 50 55 54 79 70  |7.........CPUTyp|
00000320  65 31 00 00 00 00 0d 00  00 00 01 6c 69 63 65 6e  |e1.........licen|
00000330  73 65 5f 63 6f 75 6e 74  31 00 00 00 00 0c 00 00  |se_count1.......|
00000340  00 28 36 38 65 63 37 33  33 34 31 32 62 39 41 32  |.(68ec733412b9A2|
00000350  30 35 46 37 33 42 44 41  35 35 37 46 36 36 39 44  |05F73BDA557F669D|
00000360  41 43 44 30 41 42 33 30  42 37 46 36 35 36 33 33  |ACD0AB30B7F65633|
00000370  35 39 36 32 30 41 00 00  00 00 12 00 00 00 11 70  |59620A.........p|
00000380  61 72 65 6e 74 5f 73 68  6f 72 74 5f 74 69 74 6c  |arent_short_titl|
00000390  65 59 6f 75 72 20 46 69  72 73 74 20 4c 69 73 74  |eYour First List|
000003a0  65 6e 00 00 00 00 0c 00  00 00 11 70 61 72 65 6e  |en.........paren|
000003b0  74 5f 74 69 74 6c 65 59  6f 75 72 20 46 69 72 73  |t_titleYour Firs|
000003c0  74 20 4c 69 73 74 65 6e  00 00 00 00 0e 00 00 00  |t Listen........|
000003d0  0b 70 75 62 5f 64 61 74  65 5f 73 74 61 72 74 30  |.pub_date_start0|
000003e0  31 2d 4a 41 4e 2d 32 30  30 30 00 00 00 00 11 00  |1-JAN-2000......|
000003f0  00 00 2c 73 68 6f 72 74  5f 64 65 73 63 72 69 70  |..,short_descrip|
00000400  74 69 6f 6e 59 6f 75 72  20 41 75 64 69 62 6c 65  |tionYour Audible|
00000410  20 61 64 76 65 6e 74 75  72 65 73 20 62 65 67 69  | adventures begi|
00000420  6e 20 72 69 67 68 74 20  68 65 72 65 2e 2e 2e 2e  |n right here....|
00000430  00 00 00 00 0a 00 00 00  0d 75 73 65 72 5f 61 6c  |.........user_al|
00000440  69 61 73 41 41 4a 4f 48  38 35 53 35 30 45 5a 5a  |iasAAJOH85S50EZZ|
00000450  69 0f 46 52 5f 50 52 4d  4f 5f 30 30 30 30 30 31  |i.FR_PRMO_000001|
00000460  00 00 01 e3 09 31 b9 75  92 38 17 08 62 c4 a8 0b  |.....1.u.8..b...|
00000470  da b2 be 96 ef 27 e2 2d  09 15 de 3c fa 09 3e 4e  |.....'.-...<..>N|
00000480  2d 43 93 54 a6 a6 86 79  33 03 00 12 71 70 00 00  |-C.T...y3...qp..|
00000490  01 30 00 00 00 14 00 00  00 00 00 00 00 00 00 00  |.0..............|
000004a0  00 01 a5 ab b0 fa 00 00  00 01 00 00 00 03 ff ff  |................|
....
```

Now that is interesting. We have some plain text data which seems to relate
to the book's title, a short description and there is some info about
encryption blocks and audio codecs. Let's look at this in more detail.


### Reconstructing the file header

It appears that redundancy is something Audible.com values highly. The mount
of times the book title is listed in the first 1024 bytes of the file is a bit
silly. There's probably a reason for that, but honestly, I don't care.
Let's move on to the more interesting bits.

There is repeated mention of the codec used to encode the audio data 
(MP3_32 in this case). This is useful information, as it will help us configure
the [ffmpeg][ffmpeg] decoder we'll be using when we try to process actual
audio data. If not for this, then at least it gives us valuable information
about the type of data we expect to find behind the encryption layer. The MP3
container specification is easily found on the internet.

[ffmpeg]: https://www.ffmpeg.org/

The command line tools can help us only so much. To go deeper, I would prefer
writing some code which can process the file header and extract further data
more precisely. The data we have so far gives us some hints as to the structure
of the file header.

When going back to the hex dump view, I noticed that every pair of key/value
strings is preceeded by a block of 9 bytes:

* 1 byte with unknown purpose (always `00` as far as I can tell)
* 1 32-bit integer holding the length of the key string.
* 1 32-bit integer holding the length of the value string.

```
000000b0                                       00 00 00 00  |             ...|
000000c0  0a 00 00 00 14 70 72 6f  64 75 63 74 5f 69 64 46  |.....product_idF|
000000d0  52 5f 50 52 4d 4f 5f 30  30 30 30 30 31 5f 6d 70  |R_PRMO_000001_mp|
000000e0  33 33 32                                          |332
```

We start with `00`, which I do not know the purpose of. Then `00 00 00 0a`
(decimal 10), followed by `00 00 00 14` (decimal 20). Then followed by
`product_id` (10 characters) and `FR_PRMO_000001_mp332` (20 characters).

This pattern repeats for a number of times for each key/value pair. I then
noticed the entire block of key/value pairs is preceeded with a 32-bit integer
`00 00 00 17` (decimal 23), which is awfully close to the number of pairs
we've seen so far. So this looks to be a dynamically sized dictionary, starting
with the total number of pairs and then each pair individually.

Putting all this knowledge into a simple Go program, gives us the following
output. Note that these are not in the same order as they are defined in
the file:

```
{
  Tags: {
	"68ec733412b9": "A205F73BDA557F669DACD0AB30B7F6563359620A",
	"CPUType": "1",
	"EncryptedBlocks": "52189",
	"HeaderKey": "3498980186 198870190 3001429271 2115137088",
	"HeaderSeed": "1528285495",
	"codec": "mp332",
	"copyright": "(C)Audible.com",
	"description": "Your Audible adventures begin right here....",
	"is_aggregation": "no",
	"license_count": "1",
	"license_list": "103143417",
	"long_description": "Your Audible adventures begin right here, with \"Dear Amanda\" by Steve Martin. Listen now.",
	"parent_short_title": "Your First Listen",
	"parent_title": "Your First Listen",
	"product_id": "FR_PRMO_000001_mp332",
	"provider": "Audible, Inc.",
	"pub_date_start": "01-JAN-2000",
	"pubdate": "01-JAN-2000",
	"short_description": "Your Audible adventures begin right here....",
	"short_title": "Your First Listen",
	"title": "Your First Listen",
	"title_id": "FR_PRMO_000001",
	"user_alias": "AAJOH85S50EZZ"
  }
}
```

Aces! My attention is drawn in particular to the encryption and header related
tags, along with the cryptic `68ec733412b9` tag. I'll bet that some of this is
intended to facilitate file/header checksumming. We'll leave this for later and
focus our attention on the block of 189 bytes we skipped in order to reach the
dictionary.

```
00000000  00 12 8b 6e 57 90 75 36  00 00 00 0c 00 00 00 00  |...nW.u6........|
00000010  00 00 00 00 00 00 00 00  00 12 8b 6e 00 00 00 01  |...........n....|
00000020  00 00 00 a0 00 00 00 18  00 00 00 02 00 00 00 b8  |................|
00000030  00 00 03 98 00 00 00 03  00 00 04 50 00 00 00 4a  |...........P...J|
00000040  00 00 00 04 00 00 04 9a  00 00 00 04 00 00 00 05  |................|
00000050  00 00 04 9e 00 00 00 08  00 00 00 06 00 00 04 a6  |................|
00000060  00 00 09 9e 00 00 00 07  00 00 0e 44 00 00 01 44  |...........D...D|
00000070  00 00 00 08 00 00 0f 88  00 00 00 04 00 00 00 09  |................|
00000080  00 00 0f 8c 00 00 00 00  00 00 00 0a 00 00 0f 8c  |................|
00000090  00 12 71 88 00 00 00 0b  00 12 81 14 00 00 0a 5a  |..q............Z|
000000a0  00 01 00 00 57 90 75 36  00 00 00 00 00 00 00 00  |....W.u6........|
000000b0  00 00 00 00 00 00 00 00                           |........        |
...
```

There has to be some purpose to all of this. Going by what little experience I
have with (binary) file formats, this usually contains things like a magic
number, file version and table of contents. Let's put that assumption to the
test.

The first four bytes in the file, read as a 32-bit integer, yield: `1215342`.
That looks like something I've seen before...

```
$ ls -l data/
-rw-r----- 1 jim jim 1215342 Nov  5 17:20 file.aa
-rw-r--r-- 1 jim jim    2642 Nov  6 12:39 file.jpeg
-rw-r--r-- 1 jim jim    4096 Nov  5 23:32 header.bin
```

Neat! It's exactly the file size. The next four bytes (`57 90 75 36`) read as
`1469084982`. That doesn't ring a bell. A quick search in the data shows that
this same number is repeated at offset `000000a0 + 4`. It may serve as some
kind of delimiter. Let's leave it for later.

Next is, what looks like, another 32 bit integer `00 00 00 0c` (decimal __12__).
I went ahead and assumed this indicated the size of some block of bytes
directly following it. But a little bit further into the data, the file size
presents itself again and there appears to be a repeating pattern of 3
integers, where the first integer has an incrementing value (`00 00 00 01`,
`00 00 00 02`, `00 00 00 03`, ...., `00 00 00 0b`) (decimals 0 through 11).
That's __12__ repeating blocks. This can't be a coincidence.

Examining the other two values in each 3-integer sequence, quickly reveals
that this table relates to block offsets and sizes. Where the very first entry
defines the entire file as a whole:

	00000010  00 00 00 00 00 00 00 00  00 12 8b 6e              |...........n    |

Table entry `00 00 00 00` (0), offset `00 00 00 00` (0) and size `00 12 8b 6e`,
which we already know to be the file size. The next entry only confirms our
suspicions:

	00000010                                       00 00 00 01  |            ....|
	00000020  00 00 00 a0 00 00 00 18                           |........        |

Table entry `00 00 00 01` (1), offset `00 00 00 a0` (160), size `00 00 00 18`
(24). The offset in particular is revealing. If we look at offset `000000a0`
in the header data, we see that this is exactly the place where this 12-entry
table ends.

The only part unaccounted for, is the group of 4 bytes directly following the
table size. Putting it all together in code, gives us the following list.
I have named it 'TOC' for Table of Contents, because that's what this appears
to be:

```
TOC: [
	[0, 1215342],
	[160, 24],
	[184, 920],
	[1104, 74],
	[1178, 4],
	[1182, 8],
	[1190, 2462],
	[3652, 324],
	[3976, 4],
	[3980, 0],
	[3980, 1208712],
	[1212692, 2650]
]
```

Note how these values all nicely tie together.

	toc[1]       160 +   24 =  184
	toc[2]       184 +  920 = 1104
	toc[3]      1104 +   74 = 1178
	             ...    ...    ...
	toc[11]  1212692 + 2650 = 1215342 = total file size

This is a pretty solid indicator that we're doing the right thing.
What we are left with now, is a sequence of 24 bytes directly after the table:

	000000a0  00 01 00 00 57 90 75 36  00 00 00 00 00 00 00 00  |....W.u6........|
	000000b0  00 00 00 00 00 00 00 00                           |........        |

In the previous paragraph, we found a TOC entry that points to this exact
offset and length. We still don't know what this block means. All we know is
it contains the same unidentified (magic?) number we already saw as the second
entry in the header. So let's just call this the header terminator.

Another thing to address is the values of all the TOC offsets. We would like
to understand what each of these offsets are for. There are a few we can
already fill in and some others we can make an educated guess about:

	TOC[0] = [0, 1215342]        => Defines entire file
	TOC[1] = [160, 24]           => Header termination block?
	TOC[2] = [184, 920]          => Dictionary with tags.
	
	TOC[10] = [3980, 1208712]    => Huge block; probably the audio.
	TOC[11] = [1212692, 2650]    => This offset is very close to the JPEG
	                                we extracted at the beginning. The size
	                                is also pretty much the same. This is
	                                the cover image block.

## Audio

Now that we know where the audio data is located, we can focus our attention
on getting it out. We have no idea what this file is encoded as. One thing's
for sure, it is not plain MP3 data. This is where the hard part begins.


### Digital Rubbish Management

I have zero experience in dealing with DRM schemes, so I was in dire need of
some internet wisdom. Some more Googling later, I ran into this
[PDF](http://esec-lab.sogeti.com/dotclear/public/publications/10-hitbkl-drm.pdf),
which goes into some more depth on the subject. It does not provide any
technical details, other than some hints on where to start looking.

What's left for me, is to actually go ahead and do it. I will probably start
by installing Audible's own audio manager and try to hook up a debugger to
see what it's doing when loading an audio book.

To be continued!

