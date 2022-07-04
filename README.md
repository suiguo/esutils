# esutils
convert  struct with json tag to elasticsearch mapping.

its a very simple code.u can edit it by yourself,if u want to add new type. just support elasticsearch v7 and v8.

| go type  |  es type |
| ------------- |:-------------:|
| uint8,uint16,uint32,int8,int16,int32,int| long     |
|uint64,int64    | long     |
| float32      | float    |
| float64     | double    |
| time.Time,*time.Time      | date    |
| string      |  keyword(default ignore_above 256)    |



将带有json tag的go 结构体转换成 elasticsearch的mapping.


import "github.com/suiguo/esutils" 

https://github.com/suiguo/esutils/blob/main/example/example.go
