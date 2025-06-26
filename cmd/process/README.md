# What?

Use ffmpeg for
- transcode audio file.
- fade in/out audio file.
- mix multiple audio files.

# Examples

## Transcode

bin/process -job '{"input":"bin/elefanten.mp3","outputs":["bin/output_128.mp3","bin/output_196.mp3"],"type":"transcode",audioFormats":[{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2},{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}]}' -codecConfig cmd/process/codec-config.yaml

## Fade job

### Fade in

bin/process -job '{"input":"bin/elefanten.mp3","outputs":["bin/output_128.mp3","bin/output_196.mp3"],"type":"fadeIn","procInfo":{"offset":3.0,"duration":9.0,"from":0.0,"to":1.0},"audioFormats":[{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2},{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}]}' -codecConfig cmd/process/codec-config.yaml

### Fade out

bin/process -job '{"input":"bin/elefanten.mp3","outputs":["bin/output_128.mp3","bin/output_196.mp3"],"type":"fadeIn","procInfo":{"offset":3.0,"duration":9.0,"from":0.0,"to":1.0},"audioFormats":[{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2},{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}]}' -codecConfig cmd/process/codec-config.yaml
