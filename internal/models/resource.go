package models

type Resources struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
	GPU    int `json:"gpu"`
}

func (r Resources) IsZero() bool {
	return r.CPU == 0 && r.Memory == 0 && r.GPU == 0
}

func (r Resources) Add(other Resources) Resources {
	return Resources{
		CPU:    r.CPU + other.CPU,
		Memory: r.Memory + other.Memory,
		GPU:    r.GPU + other.GPU,
	}
}

func (r Resources) Sub(other Resources) Resources {
	return Resources{
		CPU:    r.CPU - other.CPU,
		Memory: r.Memory - other.Memory,
		GPU:    r.GPU - other.GPU,
	}
}

func (r Resources) Fits(capacity Resources) bool {
	return capacity.CPU >= r.CPU && capacity.Memory >= r.Memory && capacity.GPU >= r.GPU
}
