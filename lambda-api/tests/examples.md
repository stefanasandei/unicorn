stdin, multiple files:
```json
{
  "runtime": {
    "name": "python3",
    "version": "3.12"
  },
  "project": {
    "entry": "import utils\nprint(utils.add(int(input()),2))",
    "files": [{
      "name": "utils.py",
      "contents": "def add(a, b):\n\treturn a + b"
    }]
  },	  
  "process": {
    "stdin": "2"
  }
}

```

time limit: 
```json
{
  "runtime": {
    "name": "go",
    "version": "3.12"
  },
  "project": {
    "entry": "package main\nimport \"fmt\"\n\nfunc main() {\nfmt.Println(\"works!\")\n}"
  },
  "process": {
    "time": "5ms"
  }
}
```

env vars:
```json
{
  "runtime": {
    "name": "python3", 
    "version": "3.12"
  }, 
  "project": {
    "entry": "import os\nprint(os.environ['A'])"
  },
  "process": {
    "env": {
      "A": "W",
      "C": "D"
    }
  }
}
```