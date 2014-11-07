## audible

Package audible reads header information from Audible.com audio book files.
It takes the contents of an encrypted `.aa` file and returns the tags and
other metadata in its header.

It should be noted that not all values in the header are understood. This
file format is a proprietary format and no public documentation is available.
The way I came by the information in this package, is by analysing its
contents in a hex editor and piecing the header together from there.
Contrary to the audio data, the header is not encrypted and can be readily
accessed.

If you are interested in finding out how I pieced this together, you can read
the `MAKINGOF.md` file in this repository.

The data presented in this document comes from a sample audio book.
This is the free introductory audio book you get when you sign up for an
account on Audible.com. This is also the file I used to piece all the
information together.

All data in these files is encoded as big endian. The basic layout of a
`.aa` file is as follows:


### File layout

```
   SIZE | DESCRIPTION
--------|------------------------------------------------------------------------
      4 | File size.
      4 | File magic value.
      4 | N number of entries in table of contents.
      4 | Unidentified integer: possibly marker for beginning of TOC entries.
   N*12 | N number of TOC entries. Each consisting of 3 32-bit integers:
        | * 32-bit TOC table index: 0..N-1,
        | * 32-bit Block offset in bytes, counting from the start of the file.
        | * 32-bit Block size in bytes.
      4 | Unidentified integer: possibly marker for end of TOC entries.
     20 | Header termination block; No other known purpose. Starts with file
        | magic value, followed by 16 zero bytes.
      4 | M number of entries in tag dictionary.
M*(9+?) | M number of key/value pairs, each representing a single tag.
        | * 1 unidentified byte
        | * 32-bit key string length in bytes (X).
        | * 32-bit value string length in bytes (Y).
        | * X bytes of key string data.
        | * Y bytes of value string data.
    N*? | N number of data blocks defining anything from the encrypted audio
        | to a cover image to other stuff. The size of each block is defined
        | by the respective entry in the TOC.
```

### Table of contents.

The TOC defines a list of integer tuples. Each one holds the offset and size
of an important block in the file. Here's a listing of the blocks and a note
on which ones I was able to identify:

```
 TOC[0] = [0, 1215342]        => Defines entire file
 TOC[1] = [160, 24]           => Header termination block?
 TOC[2] = [184, 920]          => Dictionary with tags.
 TOC[3] = [1104, 74]          => Unknown.
 TOC[4] = [1178, 4]           => Unknown.
 TOC[5] = [1182, 8]           => Unknown.
 TOC[6] = [1190, 2462]        => Unknown.
 TOC[7] = [3652, 324]         => Unknown.
 TOC[8] = [3976, 4]           => Unknown.
 TOC[9] = [3980, 0]           => Unknown.
TOC[10] = [3980, 1208712]     => Huge block; very likely the audio.
TOC[11] = [1212692, 2650]     => The cover image block. This contains a raw,
                                 unencrypted JPEG image in the case of our test
                                 file.
```

### Tags

There are a number of tags defined in the file I used to create this package.
From what I can gather from some cursory Googling, other files may contain
more tags, but here's the rundown of the ones I know of.

The purpose for most of these is pretty self-evident. There are a few in here
I can't determine the meaning of, but they are most likely related to the DRM
encryption and checksumming of the header/file.

```
"long_description" = "Your Audible adventures begin right here, with \"Dear Amanda\" by Steve Martin. Listen now."
"copyright" = "(C)Audible.com"
"pub_date_start" = "01-JAN-2000"
"title_id" = "FR_PRMO_000001"
"HeaderKey" = "3498980186 198870190 3001429271 2115137088"
"short_title" = "Your First Listen"
"is_aggregation" = "no"
"license_count" = "1"
"68ec733412b9" = "A205F73BDA557F669DACD0AB30B7F6563359620A"
"parent_title" = "Your First Listen"
"product_id" = "FR_PRMO_000001_mp332"
"provider" = "Audible, Inc."
"pubdate" = "01-JAN-2000"
"short_description" = "Your Audible adventures begin right here...."
"HeaderSeed" = "1528285495"
"EncryptedBlocks" = "52189"
"license_list" = "103143417"
"CPUType" = "1"
"parent_short_title" = "Your First Listen"
"title" = "Your First Listen"
"description" = "Your Audible adventures begin right here...."
"codec" = "mp332"
"user_alias" = "AAJOH85S50EZZ"
```


### Usage

    go get github.com/jteeuwen/audible


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

