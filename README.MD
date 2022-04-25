# converter

> утилита предназначена для конвертации файлов pdf в jpg с полным повторением структуры каталогов


## пример

структура папки source:
```
- test-folder1      (папка)
    - pdf-test.pdf  (файл)
- test-folder2      (папка)
    - test.pdf      (файл)  
```

структура папки target:
```
- test-folder1      (папка)
    - pdf-test.pdf  (папка)
        - 1.jpg     (файл)
- test-folder2      (папка)
    - test.pdf      (папка)   
        - 1.jpg     (файл)
        - 2.jpg     (файл)
        - 3.jpg     (файл)
        - ..jpg     (файл)
```

## использование

windows:
```
converter.exe -src "C:\folder\source" -tar "C:\folder\target"
```

linux:
```
converter.exe -src "./source" -tar "./target"
```

## сборка из исходников
```
go build -o converter ./cmd/main.go
```

## зависимости
```
github.com/karmdip-mi/go-fitz
```