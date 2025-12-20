#!/bin/bash

echo "🌸 开始运行下载器单元测试..."


echo "1. 运行基本测试..."
go test ./test -v

echo ""
echo "2. 运行测试并计算覆盖率..."
go test ./test -cover -v

echo ""
echo "3. 生成详细的覆盖率报告..."
go test ./test -coverprofile=test/coverage.out
go tool cover -html=test/coverage.out -o test/coverage.html

echo ""
echo "4. 输出覆盖率摘要..."
go tool cover -func=test/coverage.out

echo ""
echo "测试完毕！"
echo "覆盖率报告已生成: test/coverage.html"