drawio-read
============

One day I saw some magic. I saw that you could export a PNG from [draw.io](https://www.draw.io/), then any time later, import that PNG into draw.io and it would fully rebuild the editable graph. This blew my mind. How could this know layers and colors and attributes and settings.. from a raster image?

It turns out that draw.io _doesn't_ try to recreate the graph from the image data. That would be a feat, surely. Rather, it embeds the information it needs into meta-information (like EXIF) of the PNG. Curious still, I tried to use a few standard meta-info reading tools like [`exiftool`](https://www.sno.phy.queensu.ca/~phil/exiftool/) and uncompressing with `uncompress` or `zlib-flate`, but to no luck. exiftool choked on some invalid CRC bits, and the decompression tools had header problems.

So, I wrote this little parser to read one of these meta-enriched PNGs and extract that information. For funs.

To use:

```sh
./read.py path/to/image.png
```

It will then print the meta information to the terminal. 

Or save it with 

```sh
./read.py path/to/image.png > diagram.xml
```

---

If you're curious about the format:

The zTXT section of the PNG meta has key/value properties. The data is stored under 'mxGraphModel'. The data itself is compressed and base64'd XML, embedded as text in another XML document, which is compressed and stored as the contents of the 'mxGraphModel' zTXT section. There are a few peculiarities like invalid CRCs, and header-less compression, which made this custom tool necessary over exiftool and pals.

exiftool output:

```sh
$ exiftool file.png 
ExifTool Version Number         : 10.55
File Name                       : file.png
File Size                       : 11 kB
File Modification Date/Time     : 2018:02:15 17:22:33-05:00
File Access Date/Time           : 2018:02:15 17:22:33-05:00
File Inode Change Date/Time     : 2018:02:15 17:47:58-05:00
File Permissions                : rw-r--r--
File Type                       : PNG
File Type Extension             : png
MIME Type                       : image/png
Image Width                     : 371
Image Height                    : 451
Bit Depth                       : 8
Color Type                      : RGB with Alpha
Compression                     : Deflate/Inflate
Filter                          : Adaptive
Interlace                       : Noninterlaced
Warning                         : Error inflating mxGraphModel
Mx Graph Model                  : (Binary data 857 bytes, use -b option to extract)
Image Size                      : 371x451
Megapixels                      : 0.167
```

looking at the mxgraphmodel info (using xxd to view the binary info)

```sh
$ exiftool -b -mxgraphmodel file.png | xxd -ps
Warning: Error inflating mxGraphModel - file.png
c9b6aa3a10fd1ad7ba77705c904080218d200710a4912393b7682244e944
3afdfa179f67f046a9ecaabd2ba9546503d566bd901a6f00333df02097b8
1d3750db00e0742f52d7e906e8fc96a1ee3f3f2cbb810ab56cd24e2b5d57
11fd83b8bfd492fbbec631ce2c32bee3a1b085e84db1f6a1636f804aed9a
dcde490c9cdfba3745ad86aea1888eb82db38540045b967be709d24b3a90
ffc9000ace787890aefd9c4cdcb25bf183e3828cddf0819765d91643ba6c
49f7718ecf1e7f5c059e498edf28dc6da05a90b41cd2868690e213c06602
0280075f0c5b705f18e4ecd705e5f8abe091280194231ec38f689b36bfa2
5e5ae22ff623fa3a15879c7f2c3b8f617f34577e99d81bf0a3a9e0b1715e
1ce3e6f7450e23d5c9cc85d1a58b9c4f7b8e842074e16eae2b04fa136825
3465e6f478cc7dcc105bfa8eedb26ee6fd189d3640892d61b08d35bb0827
9379ba8a4e2b94c866dd387ee52ee71f51395e7c3ba832ed7555a263abe9
6ae53afe37aff5f2d309c8ed4e8babac49ee79514992453616f22ca9880f
9f6b9c9ea9d368afb02a443528b46b62f83666a6b2bea5b623862ffb3ba7
b1454fefabbbc8967fc024054637de95de60fd5312a887d87bb49d3ebb27
d05990568bb68adebaa822b0d29f05af69bb2ee9aa822b593d2c4fd13d22
bbb0add1d24a708af6c21181f61edd39387448dd47a106ddf54655c6e719
b40994e4fb319b9b24b152e34abb526191290f8ccb89382aee56b48bb2ae
235617caea18e9fc28680e7c79b7d89713eb30f6b27cc79690b3437f2771
6207997170b46f971de5ab117b38aa35fd424549ad8603e96bf57ad685c4
ea0baf39eb366f2f0aba95f3a0f166722a7b6cadf8825c4b6cac70d2e483
9b5c652cd0eb5285bd3671011694696f7c9b96dc2f142b8d23f27bad1bf4
267c61a79b80802a636cd940e6594b134e7ee7d80f97e7d726b70566e00f
94a51f281317e895f769e8bc2c39b4d918019171ec5089a169f2499e256d
31c21fa99c05133e24a4bae7ac951c3ef0959b7b3804e1dcad827f2a5ba4
922a7526f57a0d417d5c3dec48a613beab78e0aa1931f38e8c7b7b852d7b
a6a85ec5527511d8f91cd117d01d0b3adacb5dd94ad4629c1afb26af0de8
711e0457f5c907bc5b3fde2cda01b617f10a72678084d39c490a56074f54
```


trying to uncompress, getting hit with header (wbits) issues. Couldn't find any flags to set this with `zlib-flate`, `gzip`, `uncompress`, etc.

```sh
exiftool -b -mxgraphmodel file.png | zlib-flate -uncompress
Warning: Error inflating mxGraphModel - file.png
flate: inflate: data: incorrect header check
```

Using this tool:

```sh
$ ./read.py file.png 
```

```xml
<mxGraphModel dx="1400" dy="758" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="850" pageHeight="1100" background="#ffffff" math="0" shadow="0">
  <root>
    <mxCell id="0"/>
    <mxCell id="1" parent="0"/>
    <mxCell id="10" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;exitX=1;exitY=0.5;jettySize=auto;orthogonalLoop=1;" parent="1" source="2" target="8" edge="1">
      <mxGeometry relative="1" as="geometry"/>
    </mxCell>
    <mxCell id="2" value="" style="rounded=0;whiteSpace=wrap;html=1;" parent="1" vertex="1">
      <mxGeometry x="130" y="260" width="120" height="60" as="geometry"/>
    </mxCell>
    <mxCell id="5" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;exitX=0.13;exitY=0.77;exitPerimeter=0;entryX=0.5;entryY=0;jettySize=auto;orthogonalLoop=1;" parent="1" source="3" target="2" edge="1">
      <mxGeometry relative="1" as="geometry"/>
    </mxCell>
    <mxCell id="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;exitX=0.55;exitY=0.95;exitPerimeter=0;entryX=0.5;entryY=0;jettySize=auto;orthogonalLoop=1;" parent="1" source="3" target="6" edge="1">
      <mxGeometry relative="1" as="geometry">
        <Array as="points">
          <mxPoint x="431" y="205"/>
          <mxPoint x="440" y="205"/>
        </Array>
      </mxGeometry>
    </mxCell>
    <mxCell id="3" value="" style="ellipse;shape=cloud;whiteSpace=wrap;html=1;" parent="1" vertex="1">
      <mxGeometry x="365" y="70" width="120" height="80" as="geometry"/>
    </mxCell>
    <mxCell id="9" style="edgeStyle=orthogonalEdgeStyle;rounded=0;html=1;exitX=0;exitY=0.5;entryX=0.5;entryY=0;jettySize=auto;orthogonalLoop=1;" parent="1" source="6" target="8" edge="1">
      <mxGeometry relative="1" as="geometry"/>
    </mxCell>
    <mxCell id="6" value="" style="rounded=0;whiteSpace=wrap;html=1;" parent="1" vertex="1">
      <mxGeometry x="380" y="260" width="120" height="60" as="geometry"/>
    </mxCell>
    <mxCell id="8" value="" style="shape=cylinder;whiteSpace=wrap;html=1;boundedLbl=1;" parent="1" vertex="1">
      <mxGeometry x="290" y="440" width="60" height="80" as="geometry"/>
    </mxCell>
  </root>
</mxGraphModel>
```