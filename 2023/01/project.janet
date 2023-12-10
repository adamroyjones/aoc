(declare-project
 :name "day-01"
 :description "Attempts at solutions for day 1"
 :dependencies [{:url "https://github.com/ianthehenry/judge.git" :tag "v2.7.2"}])

(task "test" [] (shell "jpm_tree/bin/judge"))
