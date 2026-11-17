# OBSCURE

OBSCURE will be a command-line tool for selectively redacting video.

```sh
make
obscure video.mp4
```

The input must contain a video stream. FFmpeg decodes it to packed RGB24 frames. The pipeline passes the frames to the a backend.
