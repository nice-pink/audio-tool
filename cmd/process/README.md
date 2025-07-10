# What?

Use ffmpeg for
- transcode audio file.
- fade in/out audio file.
- mix multiple audio files.

# Examples

## Transcode

bin/process -job '{"type":"transcode","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128.mp3","format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}]}' -codecConfig cmd/process/codec-config.yaml

## Transcode - Segment

bin/process -job '{"type":"transcode","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128_%03d.mp3","segmentDuration":5.0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196_%03d.mp3","segmentDuration":5.0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}]}' -codecConfig cmd/process/codec-config.yaml

## Fade job

### Fade in

bin/process -job '{"type":"fade","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}],"procInfo":{"offset":3.0,"duration":9.0,"from":0.0,"to":1.0}}' -codecConfig cmd/process/codec-config.yaml

### Fade in - Trim start

bin/process -job '{"type":"fade","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}],"procInfo":{"offset":3.0,"duration":9.0,"from":0.0,"to":1.0,"trim":true}}' -codecConfig cmd/process/codec-config.yaml

### Fade out

bin/process -job '{"type":"fade","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}],"procInfo":{"offset":3.0,"duration":9.0,"from":0.0,"to":1.0}}' -codecConfig cmd/process/codec-config.yaml

### Fade out - Trim end

bin/process -job '{"type":"fade","input":"bin/elefanten.mp3","outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}],"procInfo":{"offset":3.0,"duration":9.0,"from":1.0,"to":0.0,"trim":true}}' -codecConfig cmd/process/codec-config.yaml

## Mix

bin/process -job '{"type":"mix","inputs":[{"filename":"bin/elefanten.mp3"},{"filename":"bin/elefanten.mp3","offset":5.0}],"procJob":{"outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}]}}' -codecConfig cmd/process/codec-config.yaml

Segments:

bin/process -job '{"type":"mix","inputs":[{"filename":"bin/elefanten.mp3"},{"filename":"bin/elefanten.mp3","offset":5.0}],"procJob":{"outputs":[{"filename":"bin/output_128_%03d.mp3","segmentDuration":5,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196_%03d.mp3","segmentDuration":5,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}]}}' -codecConfig cmd/process/codec-config.yaml

Single output:

bin/process -job '{"type":"mix","inputs":[{"filename":"bin/elefanten.mp3"},{"filename":"bin/elefanten.mp3","offset":5.0}],"procJob":{"outputs":[{"filename":"bin/output_128.mp3","segmentDuration":0,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}}}]}}' -codecConfig cmd/process/codec-config.yaml

With processing:

bin/process -job '{"type":"mix","inputs":[{"filename":"bin/elefanten.mp3"},{"filename":"bin/elefanten.mp3","offset":5.0}],"procJob":{"outputs":[{"filename":"bin/output_128.mp3","segmentDuration":5,"format":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}},{"filename":"bin/output_196.mp3","segmentDuration":5,"format":{"codec":"mp3","bitrate":196000,"sampleRate":44100,"channels":2}}]},"procInfos":[{"offset":0.0,"duration":3.0,"from":1.0,"to":0.0,"trim":true},{"offset":3.0,"duration":3.0,"from":1.0,"to":0.0,"trim":true}]}' -codecConfig cmd/process/codec-config.yaml
