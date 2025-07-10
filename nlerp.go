package nonlinear

/*
 * Non-linear interpolations between 0 and 1.
 * Clamping is enforced lest the result not be defined outside of [0,1].
 */

// NLerp returns the value of the supplied non-linear function at t. Note t is clamped to [0,1]
func NLerp(t, start, end float64, f NonLinear) float64 {
	if t < 0 {
		return start
	}
	if t > 1 {
		return end
	}
	t = f.Transform(t)
	return (1-t)*start + t*end
}

// InvNLerp performs the inverse of NLerp and returns the value of t for a value v (clamped to [start, end]).
func InvNLerp(v, start, end float64, f NonLinear) float64 {
	t := (v - start) / (end - start)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return f.InvTransform(t)
}

// RemapNL converts v from one space to another by applying InvNLerp to find t in the initial range, and
// then using t to find v' in the new range.
func RemapNL(v, istart, iend, ostart, oend float64, fi, fo NonLinear) float64 {
	return NLerp(InvNLerp(v, istart, iend, fi), ostart, oend, fo)
}
