package metrics

import (
	"bytes"
	"encoding/binary"
	"errors"
	"nms/internal/utils"
)

type Metrics struct {
	PacketID byte
	AgentID  byte
	TaskID   uint16
	Time     string
	Metrics  string
}

type MetricsBuilder struct {
	Metrics Metrics
}

func NewMetricsBuilder() *MetricsBuilder {
	return &MetricsBuilder{
		Metrics: Metrics{
			PacketID: 0,
			AgentID:  0,
			TaskID:   0,
			Time:     "",
			Metrics:  "",
		},
	}
}

func (m *MetricsBuilder) SetPacketID(packetID byte) *MetricsBuilder {
	m.Metrics.PacketID = packetID
	return m
}

func (m *MetricsBuilder) SetAgentID(agentID byte) *MetricsBuilder {
	m.Metrics.AgentID = agentID
	return m
}

func (m *MetricsBuilder) SetTaskID(taskID uint16) *MetricsBuilder {
	m.Metrics.TaskID = taskID
	return m
}

func (m *MetricsBuilder) SetTime(time string) *MetricsBuilder {
	m.Metrics.Time = time
	return m
}

func (m *MetricsBuilder) SetMetrics(metrics string) *MetricsBuilder {
	m.Metrics.Metrics = metrics
	return m
}

func (m *MetricsBuilder) Build() Metrics {
	return m.Metrics
}

func DecodeMetrics(packet []byte) (Metrics, error) {
	buf := bytes.NewReader(packet)
	var metrics Metrics

	if len(packet) < 10 { // Adjusted length check
		return metrics, errors.New("invalid packet length")
	}

	packetID, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	agentID, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	if err := binary.Read(buf, binary.BigEndian, &metrics.TaskID); err != nil {
		return metrics, err
	}
	metrics.PacketID = packetID
	metrics.AgentID = agentID

	var timeLen uint16
	if err := binary.Read(buf, binary.BigEndian, &timeLen); err != nil {
		return metrics, err
	}

	timeBytes := make([]byte, timeLen)
	if _, err := buf.Read(timeBytes); err != nil {
		return metrics, err
	}
	metrics.Time = string(timeBytes)

	var metricsLen uint16
	if err := binary.Read(buf, binary.BigEndian, &metricsLen); err != nil {
		return metrics, err
	}

	metricsBytes := make([]byte, metricsLen)
	if _, err := buf.Read(metricsBytes); err != nil {
		return metrics, err
	}
	metrics.Metrics = string(metricsBytes)

	return metrics, nil
}

func EncodeMetrics(metrics Metrics) []byte {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(utils.METRICSGATHERING))
	buf.WriteByte(metrics.PacketID)
	buf.WriteByte(metrics.AgentID)
	binary.Write(buf, binary.BigEndian, metrics.TaskID)

	timeBytes := []byte(metrics.Time)
	timeLen := uint16(len(timeBytes))
	binary.Write(buf, binary.BigEndian, timeLen)
	buf.Write(timeBytes)

	metricsBytes := []byte(metrics.Metrics)
	metricsLen := uint16(len(metricsBytes))
	binary.Write(buf, binary.BigEndian, metricsLen)
	buf.Write(metricsBytes)

	return buf.Bytes()
}
