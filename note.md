
## 关于  val.xxxVal 、types.xxxKind
go 的内存布局关于内存申请部分, 如果一个结构体中的一个字段逃逸到堆中, 整个结构体都会逃逸到堆中
所以为了, 简化类型签名, 大量使用了类型 c 中经常使用的类似 gc_header 的方式强转了类型
转回来时候必须判断类型, 否则是个高风险的操作

