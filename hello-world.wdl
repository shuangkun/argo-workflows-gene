# Comments are allowed before versions

version 1.0

# This is how you would
# write a long
# multiline
# comment

task test {
  #This comment will not be included within the command
  command <<<
    #This comment WILL be included within the command after it has been parsed
    echo 'Hello World'
  >>>

  output {
    String result = read_string(stdout())
  }

  runtime {
    container: "my_image:latest"
  }
}

workflow wf {
  input {
    Int number  #This comment comes after a variable declaration
  }

  #You can have comments anywhere in the workflow
  call test

  output { #You can also put comments after braces
    String result = test.result
  }
}
