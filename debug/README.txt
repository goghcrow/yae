Power Debug

参考 https://groovy-lang.org/semantics.html#_power_assertion

Groovy 实现方式是在编译期间改写代码，把单行断言表达式展开成最小颗粒度，收集求值结果
这里的实现是 hook 了编译器生成的闭包, 把需要展示断言结果的 term 的求值结果记录下来

lexer parser 需要计算 loc, desugar 需要保留 loc, 部分 term (subscrip|member|call)
额外加入了 opLoc, 用来精确定位 term 求值渲染列

已知问题: 非等宽字体、非 ASCII 名称 identifier 没法对齐

渲染代码参考 groovy 编译器
