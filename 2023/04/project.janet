(declare-project
 :name "day-04"
 :description "Attempts at solutions for day 4"
 :dependencies [{:url "https://github.com/ianthehenry/judge.git" :tag "v2.8.0"}])

(task "test" [] (shell "jpm_tree/bin/judge"))
