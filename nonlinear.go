package nonlinear

import "math"

// NonLinear interface defines the transform and its inverse as used by NLerp etc.
// For mapping 0 -> 1 non-linearly. No checks! Only valid in range [0,1]
type NonLinear interface {
	Transform(t float64) float64
	InvTransform(v float64) float64
}

// NLLinear v = t
type NLLinear struct{}

func (nl *NLLinear) Transform(t float64) float64 {
	return t
}

func (nl *NLLinear) InvTransform(v float64) float64 {
	return v
}

// NLSquare v = t^2
type NLSquare struct{}

func (nl *NLSquare) Transform(t float64) float64 {
	return t * t
}

func (nl *NLSquare) InvTransform(v float64) float64 {
	return math.Sqrt(v)
}

// NLCube v = t^3
type NLCube struct{}

func (nl *NLCube) Transform(t float64) float64 {
	return t * t * t
}

func (nl *NLCube) InvTransform(v float64) float64 {
	return math.Pow(v, 1/3.0)
}

// NLExponential v = (exp(t*k) - 1) * scale
type NLExponential struct {
	K     float64
	Scale float64
}

func NewNLExponential(k float64) *NLExponential {
	return &NLExponential{k, 1 / (math.Exp(k) - 1)}
}

func (nl *NLExponential) Transform(t float64) float64 {
	return (math.Exp(t*nl.K) - 1) * nl.Scale
}

func (nl *NLExponential) InvTransform(v float64) float64 {
	return math.Log1p(v/nl.Scale) / nl.K
}

// NLLogarithmic v = log(1+t*k) * scale
type NLLogarithmic struct {
	K     float64
	Scale float64
}

func NewNLLogarithmic(k float64) *NLLogarithmic {
	return &NLLogarithmic{k, 1 / math.Log1p(k)}
}

func (nl *NLLogarithmic) Transform(t float64) float64 {
	return math.Log1p(t*nl.K) * nl.Scale
}

func (nl *NLLogarithmic) InvTransform(v float64) float64 {
	return (math.Exp(v/nl.Scale) - 1) / nl.K
}

// NLSin v = sin(t) with t mapped to [-Pi/2,Pi/2]
type NLSin struct{} // first derivative 0 at t=0,1

func (nl *NLSin) Transform(t float64) float64 {
	return (math.Sin((t-0.5)*math.Pi) + 1) / 2
}

func (nl *NLSin) InvTransform(v float64) float64 {
	return math.Asin((v*2)-1)/math.Pi + 0.5
}

// NLSin1 v = sin(t) with t mapped to [0,Pi/2]
type NLSin1 struct{} // first derivative 0 at t=1

func (nl *NLSin1) Transform(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

func (nl *NLSin1) InvTransform(v float64) float64 {
	return math.Asin(v) / math.Pi * 2
}

// NLSin2 v = sin(t) with t mapped to [-Pi/2,0]
type NLSin2 struct{} // first derivative 0 at t=0,1

func (nl *NLSin2) Transform(t float64) float64 {
	return math.Sin((t-1)*math.Pi/2) + 1
}

func (nl *NLSin2) InvTransform(v float64) float64 {
	return math.Asin(v-1)*2/math.Pi + 1
}

// NLCircle1 v = 1 - sqrt(1-t^2)
type NLCircle1 struct{}

func (nl *NLCircle1) Transform(t float64) float64 {
	if t < 1 {
		return 1 - math.Sqrt(1-t*t)
	}
	return 1
}

func (nl *NLCircle1) InvTransform(v float64) float64 {
	if v < 1 {
		return math.Sqrt(1 - (v-1)*(v-1))
	}
	return 1
}

// NLCircle2 v = sqrt(2t-t^2)
type NLCircle2 struct{}

func (nl *NLCircle2) Transform(t float64) float64 {
	return math.Sqrt(t * (2 - t))
}

func (nl *NLCircle2) InvTransform(v float64) float64 {
	return 1 - math.Sqrt(1-v*v)
}

// NLLame (aka superellipse) v = 1 - (1-t^n)^1/m
type NLLame struct {
	N   float64
	M   float64
	Odn float64
	Odm float64
}

func NewNLLame(n, m float64) *NLLame {
	return &NLLame{n, m, 1 / n, 1 / m}
}

func (nl *NLLame) Transform(t float64) float64 {
	if t < 1 {
		vm := 1 - math.Pow(t, nl.N)
		return 1 - math.Pow(vm, nl.Odm)
	}
	return 1
}

func (nl *NLLame) InvTransform(v float64) float64 {
	if v < 1 {
		v = 1 - v
		tn := 1 - math.Pow(v, nl.M)
		return math.Pow(tn, nl.Odn)
	}
	return 1
}

// NLCatenary v = cosh(t)
type NLCatenary struct{}

func (nl *NLCatenary) Transform(t float64) float64 {
	return (math.Cosh(t) - 1) / (math.Cosh(1) - 1)
}

func (nl *NLCatenary) InvTransform(v float64) float64 {
	return math.Acosh(v*(math.Cosh(1)-1) + 1)
}

// NLGauss v = gauss(t, k)
type NLGauss struct {
	K, Offs, Scale float64
}

func NewNLGauss(k float64) *NLGauss {
	offs := math.Exp(-k * k * 0.5)
	scale := 1 / (1 - offs)
	return &NLGauss{k, offs, scale}
}

func (nl *NLGauss) Transform(t float64) float64 {
	x := nl.K * (t - 1)
	x *= -0.5 * x
	return (math.Exp(x) - nl.Offs) * nl.Scale
}

func (nl *NLGauss) InvTransform(v float64) float64 {
	v /= nl.Scale
	v += nl.Offs
	v = math.Log(v)
	v *= -2
	v = math.Sqrt(v)
	return 1 - v/nl.K
}

// NLLogistic v = logistic(t, k, mp)
type NLLogistic struct {
	K, Mp, Offs, Scale float64
}

// k > 0 and mp (0,1) - not checked
func NewNLLogistic(k, mp float64) *NLLogistic {
	v0 := -mp * k
	v0 = logisticTransform(v0)
	v1 := (1 - mp) * k
	v1 = logisticTransform(v1)
	return &NLLogistic{k, mp, v0, 1 / (v1 - v0)}
}

func (nl *NLLogistic) Transform(t float64) float64 {
	t = (t - nl.Mp) * nl.K
	return (logisticTransform(t) - nl.Offs) * nl.Scale
}

func (nl *NLLogistic) InvTransform(v float64) float64 {
	v /= nl.Scale
	v += nl.Offs
	v = logisticInvTransform(v)
	return v/nl.K + nl.Mp
}

// L = 1, k = 1, mp = 0
func logisticTransform(t float64) float64 {
	return 1 / (1 + math.Exp(-t))
}

// L = 1, k = 1, mp = 0
func logisticInvTransform(v float64) float64 {
	return -math.Log(1/v - 1)
}

// NLP3 v = t^2 * (3-2t)
type NLP3 struct{} // first derivative 0 at t=0,1

func (nl *NLP3) Transform(t float64) float64 {
	return t * t * (3 - 2*t)
}

func (nl *NLP3) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

// NLP5 v = t^3 * (t*(6t-15) + 10)
type NLP5 struct{} // first and second derivatives 0 at t=0,1

func (nl *NLP5) Transform(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func (nl *NLP5) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

// NLCompound v = nl[0](nl[1](nl[2](...nl[n-1](t))))
type NLCompound struct {
	Fs []NonLinear
}

func NewNLCompound(fs []NonLinear) *NLCompound {
	return &NLCompound{fs}
}

func (nl *NLCompound) Transform(t float64) float64 {
	for _, f := range nl.Fs {
		t = f.Transform(t)
	}

	return t
}

func (nl *NLCompound) InvTransform(v float64) float64 {
	for i := len(nl.Fs) - 1; i > -1; i-- {
		v = nl.Fs[i].InvTransform(v)
	}
	return v
}

// NLOmt v = 1-f(1-t)
type NLOmt struct {
	F NonLinear
}

func NewNLOmt(f NonLinear) *NLOmt {
	return &NLOmt{f}
}

func (nl *NLOmt) Transform(t float64) float64 {
	t = 1 - t
	if t > 0 {
		return 1 - nl.F.Transform(t)
	}
	return 1
}

func (nl *NLOmt) InvTransform(v float64) float64 {
	v = 1 - v
	if v > 0 {
		return 1 - nl.F.InvTransform(1-v)
	}
	return 1
}

// NewStoppedNL uses linear interpolation between the supplied stops
type NLStopped struct {
	Stops [][]float64 // Pairs of t, v - both strictly ascending in [0,1]
}

func NewNLStopped(stops [][]float64) *NLStopped {
	// Assumes valid stops
	return &NLStopped{stops}
}

func (nl *NLStopped) Transform(t float64) float64 {
	t0, v0 := 0.0, 0.0
	ns := len(nl.Stops)
	var i int
	for i = 0; i < ns; i++ {
		if nl.Stops[i][0] > t {
			if i > 1 {
				t0 = nl.Stops[i-1][0]
				v0 = nl.Stops[i-1][1]
			}
			break
		}
	}
	if i == ns {
		t0 = nl.Stops[ns-1][0]
		v0 = nl.Stops[ns-1][1]
	}
	t1, v1 := 1.0, 1.0
	if i < ns {
		t1 = nl.Stops[i][0]
		v1 = nl.Stops[i][1]
	}
	dt := t1 - t0
	t = (t - t0) / dt
	return (1-t)*v0 + t*v1
}

func (nl *NLStopped) InvTransform(v float64) float64 {
	return bsInv(v, nl)
}

// Numerical method to find inverse
func bsInv(v float64, f NonLinear) float64 {
	n := 16
	t := 0.5
	s := 0.25

	for ; n > 0; n-- {
		if f.Transform(t) > v {
			t -= s
		} else {
			t += s
		}
		s /= 2
	}
	return t
}
