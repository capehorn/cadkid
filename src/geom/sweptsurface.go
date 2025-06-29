package geom

type SweptSurface struct {
	profile         ParametricCurve
	profileFrame    Frame
	trajectory      ParametricCurve
	trajectoryFrame Frame
}

//func (s SweptSurface) Generate(stepBy float64) Mesh {
//	trajectory := s.trajectory
//	profile := s.profile
//	tangent := trajectory.TangentAt(0, 0.01)
//	if trajectory.IsPlaneCurve {
//
//	}
//	frame := Frame{trajectory.PointAt(0), tangent}
//
//}
