lz4
===

[LZ4](https://code.google.com/p/lz4/) is a very fast lossless compression algorithm, providing compression speed at 400 MB/s per core, scalable with multi-cores CPU. It also features an extremely fast decoder, with speed in multiple GB/s per core, typically reaching RAM speed limits on multi-core systems.

This is LZ4 de/compresor for GO language. Currently supports basic LZ4 compression/decompression (no HC), and streaming format.

Speed on Intel G1610:<br>
Compresion: 209 MB/s, about 60% of C version speed<br>
Decompresion: 360 MB/s, about 20% of C version speed

If you'd like to run the benchmark through "go test -bench ." please download the Silesia corpus from [here](http://sun.aei.polsl.pl/~sdeor/corpus/silesia.zip). Unpack the files and concatinate them into one file named combined.bin.
