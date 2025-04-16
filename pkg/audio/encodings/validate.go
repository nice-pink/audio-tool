package encodings

import (
	"errors"

	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
)

// private bit validator

type PrivateBitValidator struct {
	active               bool
	verbose              bool
	metrics              *Metrics
	audioType            AudioType
	lastFrameDistance    uint64
	currentFrameCount    uint64
	parser               *Parser
	foundInitialDistance bool
}

func NewPrivateBitValidator(active bool, audioType AudioType, metrics util.MetricsControl, verbose bool) *PrivateBitValidator {
	p := &PrivateBitValidator{
		verbose:   verbose,
		active:    active,
		audioType: audioType,
		parser:    NewParser(),
	}

	if metrics.Enabled {
		p.metrics = NewMetrics(metrics.Prefix, metrics.Labels)
	}

	return p
}

func (v *PrivateBitValidator) Validate(data []byte) error {
	// bypass?
	if !v.active {
		return nil
	}

	// validate
	blockAudioInfo, err := v.parser.ParseBlockwise(data, v.audioType, true, v.verbose, false)
	if err != nil {
		log.Err(err, "Parsing error.")
		if v.metrics != nil {
			v.metrics.audioParseErrorMetric.Inc()
		}
		return err
	}

	if blockAudioInfo == nil {
		// log.Error("No block audio data.")
		return nil
	}

	// validate encodings
	for i, unit := range blockAudioInfo.Units {
		if !unit.IsPrivate {
			v.currentFrameCount++
			continue
		} else if v.verbose {
			log.Info("Found private bit.", i, v.currentFrameCount)
		}

		// validate distance
		if v.lastFrameDistance > 0 {
			if v.currentFrameCount != v.lastFrameDistance {
				if !v.foundInitialDistance {
					// skip first distance as it will usually be an error!
					v.foundInitialDistance = true
					log.Info("Skip initial distance.", v.currentFrameCount)
				} else {
					log.Error("Distances not equal. Current:", v.currentFrameCount, "!= Last:", v.lastFrameDistance, i, len(blockAudioInfo.Units))
					// update metric
					if v.metrics != nil {
						v.metrics.validationErrorMetric.Inc()
					}
				}
			} // else {
			// 	log.Info("Distance between private bits:", v.lastFrameDistance)
			// }
		}

		// reset
		v.lastFrameDistance = v.currentFrameCount
		v.currentFrameCount = 0
	}

	return nil
}

// encoding validator

type EncodingValidator struct {
	active       bool
	failEarly    bool
	expectations Expectations
	verbose      bool
	metrics      *Metrics
	parser       *Parser
}

func NewEncodingValidator(active, failEarly bool, expectations Expectations, metrics util.MetricsControl, verbose bool) *EncodingValidator {
	e := &EncodingValidator{
		expectations: expectations,
		verbose:      verbose,
		active:       active, parser: NewParser(),
	}

	if metrics.Enabled {
		e.metrics = NewMetrics(metrics.Prefix, metrics.Labels)
	}

	return e
}

func (v *EncodingValidator) Validate(data []byte) error {
	// bypass?
	if !v.active {
		return nil
	}

	// validate
	blockAudioInfo, err := v.parser.ParseBlockwise(data, GetAudioTypeFromCodecName(v.expectations.Encoding.CodecName), true, v.verbose, false)
	if err != nil {
		log.Err(err, "Parsing error.")
		return err
	}

	if blockAudioInfo == nil {
		log.Error("No block audio data.")
		return nil
	}

	// validate audio info
	isValid := IsValid(v.expectations, *blockAudioInfo, v.metrics)
	if !isValid && v.failEarly {
		return errors.New("validation failed")
	}

	// validate encodings
	for _, unit := range blockAudioInfo.Units {
		isValid = IsValidEncoding(v.expectations, unit.Encoding, v.metrics)
		if !isValid && v.failEarly {
			return errors.New("validation failed")
		}
	}

	return nil
}

// general valiation

func IsValid(expectations Expectations, audioInfo AudioInfos, metrics *Metrics) bool {
	isValid := true
	if expectations.IsCBR {
		if expectations.IsCBR != audioInfo.IsCBR {
			log.Error("IsCBR not equal:", expectations.IsCBR, "!=", audioInfo.IsCBR)
			isValid = false
		}
	}

	// update metric
	if !isValid {
		if metrics != nil {
			metrics.validationErrorMetric.Inc()
		}
	}

	return isValid
}

func IsValidEncoding(expectations Expectations, encoding Encoding, metrics *Metrics) bool {
	isValid := true
	if expectations.Encoding.Bitrate > 0 {
		if expectations.Encoding.Bitrate != encoding.Bitrate {
			log.Error("Bitrate not equal:", expectations.Encoding.Bitrate, "!=", encoding.Bitrate)
			isValid = false
		}
	}
	if expectations.Encoding.SampleRate > 0 {
		if expectations.Encoding.SampleRate != encoding.SampleRate {
			log.Error("SampleRate not equal:", expectations.Encoding.SampleRate, "!=", encoding.SampleRate)
			isValid = false
		}
	}
	if expectations.Encoding.FrameSize > 0 {
		if expectations.Encoding.FrameSize != encoding.FrameSize {
			log.Error("FrameSize not equal:", expectations.Encoding.FrameSize, "!=", encoding.FrameSize)
			isValid = false
		}
	}
	if expectations.Encoding.CodecName != "" {
		if expectations.Encoding.CodecName != encoding.CodecName {
			log.Error("CodecName not equal:", expectations.Encoding.CodecName, "!=", encoding.CodecName)
			isValid = false
		}
	}
	if expectations.Encoding.ContainerName != "" {
		if expectations.Encoding.ContainerName != encoding.ContainerName {
			log.Error("ContainerName not equal:", expectations.Encoding.ContainerName, "!=", encoding.ContainerName)
			isValid = false
		}
	}

	if expectations.Encoding.IsStereo != encoding.IsStereo {
		log.Error("IsStereo not equal:", expectations.Encoding.IsStereo, "!=", encoding.IsStereo)
		isValid = false
	}

	// update metric
	if !isValid {
		if metrics != nil {
			metrics.validationErrorMetric.Inc()
		}
	}

	return isValid
}
