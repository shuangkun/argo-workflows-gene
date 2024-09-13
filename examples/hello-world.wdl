# Comments are allowed before versions

version 1.0

# This is how you would
# write a long
# multiline
# comment

workflow wf {
  input {
    Int number  #This comment comes after a variable declaration
  }

  #You can have comments anywhere in the workflow
  call hello
  output { #You can also put comments after braces
    String result = hello.result
  }
}

task hello {
  #This comment will not be included within the command
  command <<<
    echo 'Hello World'
  >>>

  output {
    String result = read_string(stdout())
  }

  runtime {
    container: "alpine:latest"
  }
}
