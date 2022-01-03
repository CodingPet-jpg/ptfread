# ptfread

当前为per thread file reader v0.0.2，较上版本新增了json report逆序列化的功能，可以实现跨时段的增量解析



Feature:

- 1 多表唯一指定Sheet的复数列相似度比较，通过控制台输出结果与耗时信息
- 2 比较结果生成json格式报告
- 3 读取json格式报告和新的工作目录中多表进行匹配
- 4 不指定工作目录仅进行多份报告之间的比较

Config:

```yml
ActiveSheet : Sheet1
WorkDirectory : C:\Users\JOKOI\GolandProjects\testdata
Length : 5
BitMap : 8191
ParallelNum : 3000
InheritSource :
  - C:\Users\JOKOI\GolandProjects\testdata\report\report[2022-01-03 ※ 18-21-38].txt
  - C:\Users\JOKOI\GolandProjects\testdata\report\report[2022-01-03 ※ 17-32-21].txt
  - C:\Users\JOKOI\GolandProjects\testdata\report\report[2022-01-03 ※ 20-08-14].txt
```

- WorkDirectory

​		指定工作目录，即数据表所在根目录，默认匹配工作目录及其所有子目录，在工作目录下生成report目录存放比较报告

- ActiveSheet

  指定待比较的Sheet，如待比较数据表不存在此Sheet则跳过该表继续执行，并打印错误信息

- Length

​		该工具通过Length配置决定指定Sheet中的可处理行，如果当前行的有效数据列小于该配置则该行被跳过，缺省值为4

| A    | B    | C    | D    | E    | F    | G    | H    | I    |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| xxx  |      |      |      |      |      |      |      |      |
|      |      | 1    | xxx  | xxx  | xxx  | xxx  | xxx  | xxx  |
|      |      | 2    |      |      |      |      | xxx  |      |

如上表中第一行将被跳过，第二三行将被解析，有效数据列为A-I列

- BitMap

​		指定期望比较列，如上表中期望比较FGHI列则需设置数字 111100000(二进制)=》480

​		对应规则为从A列开始对应列如为期望比较列则置1，否则为0，以此类推，A列对应二进制低位，将解析完的结果转为10进制

- ParallelNum

  设定并发数，如解析途中出现异常可以调低并发数

- InheritSource

​		指定报告位置，指定复数报告不指定工作目录时将进行报告之间的比较

​		指定复数报告同时指定工作目录会读取其中一份报告作为初始值，之后将剩余报告以及工作目录中所有数据表并发地与其进行比较

​		指定工作目录即单份报告时，读取报告作为初始值，将工作目录中所有数据表相互比较后再与报告比较生成总报告

​		仅比较工作目录下数据表时请将此配置内容清空
