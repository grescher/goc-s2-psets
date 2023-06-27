package main

type User struct {
	Name   string
	Age    int
	Active bool
	Mass   float64
}

func Users() []User {
	return []User{
		{
			"John Doe",
			30,
			true,
			80.0,
		},
		{
			"Jake Doe",
			20,
			false,
			60.0,
		},
		{
			" Jane Doe ",
			150,
			true,
			.75,
		},
		{
			"\t",
			-10,
			true,
			8000.0,
		},
		{
			"\n",
			-10,
			true,
			8000.0,
		},
		{
			"Vm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=\nVm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=",
			0,
			true,
			0,
		},
		{
			"\x00\x10\x20\x30\x40\x50\x60\x70",
			0,
			true,
			0,
		},
		{
			"Billy Bones",
			-130000,
			false,
			3141567.98765456789,
		},
	}
}
