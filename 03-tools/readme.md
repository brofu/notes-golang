## Tools which can improve the productivity with Golang

#### Project Construction

* **cookiecutter**
    * Construct the project construct 
    * Can customize the template 
    * https://github.com/cookiecutter/cookiecutter
    * An template option: https://github.com/lacion/cookiecutter-golang 


#### Testing
* **gotests**
    * Insert unit test codes based on the templates
        * https://github.com/cweill/gotests
    * vim-plugin available: https://github.com/buoto/gotests-vim
         ```
        // install
        go install github.com/cweill/gotests/gotests@latest
        // install plugin
         ```
    * Need to define the templates. 
        * An example: https://github.com/ras0q/gotests-template/blob/main/templates2/function.tmpl
    * Usage
        ```
        :GoTests
        :GoTestsAll
        ```
