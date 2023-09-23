package employee

type Employee interface {
    GetPosition() string
    SetPosition(position string)
    GetSalary() float64
    SetSalary(salary float64)
    GetAddress() string
    SetAddress(address string)
}

type EmployeeStruct struct {
    position string
    salary   float64
    address  string
}

func (e *EmployeeStruct) GetPosition() string {
    return e.position
}

func (e *EmployeeStruct) SetPosition(position string) {
    e.position = position
}

func (e *EmployeeStruct) GetSalary() float64 {
    return e.salary
}

func (e *EmployeeStruct) SetSalary(salary float64) {
    e.salary = salary
}

func (e *EmployeeStruct) GetAddress() string {
    return e.address
}

func (e *EmployeeStruct) SetAddress(address string) {
    e.address = address
}