(use judge)

########################################
# General functions
########################################
(defn zip-with [f xs ys]
  (def n (min (length xs) (length ys)))
  (reduce (fn [acc el] (array/push acc (f (xs el) (ys el)))) @[] (range n)))

(defn dims [lines]
  (def nrows (length lines))
  (def column-sizes (->> lines (map length) (distinct)))
  (assert (<= (length column-sizes) 1) "the column sizes are non-uniform")
  [nrows (or (column-sizes 0) 0)])

(defn transpose [lines]
  (def in-dims (dims lines))
  (def columns (range (length (lines 0))))
  (defn construct-column [column] (map |($ column) lines))
  (def result (->> columns
                   (map construct-column)
                   (map (fn [row] (string/from-bytes ;row)))))
  (def out-dims (dims result))
  (assert (= [;(reverse in-dims)] out-dims) (printf "not a transposition (dims: %q -> %q)" in-dims out-dims))
  result)

########################################
# Functions for producing marked regions
########################################
(defn mark-special/init
  "This function replaces special characters with x and everything else with a dot.
A later function will replace each x with a 3x3 grid of xs."
  [lines]
  (def pattern ~{:replace-digit (/ (<- :d) ".")
                 :replace-special (/ (<- 1) "x")
                 :main (any (+ :replace-digit (<- ".") :replace-special))})
  (->> lines (map |(peg/match pattern $)) (map string/join)))

(defn mark-special/expand-row
  "This function replaces replaces x characters with a 1x3 vector."
  [lines]
  (def in-dims (dims lines))
  (def ncols (in-dims 1))

  (def pattern '{:dot-step (if (+ ".." (* "." -1)) (<- "."))
                 :any-x-or-x-end (if (+ ".x" "xx" (* "x" -1)) (* (constant "x") 1))
                 :x-dot (if "x." (* (constant "xx") 2)) # x.. | x.x | x.$ are all handled.
                 :main (any (+ :dot-step :any-x-or-x-end :x-dot))})

  (defn mapper [row]
    (def expanded-row (string/join (peg/match pattern row)))
    (assert (= (length expanded-row) (length row)) "the row-expansion is not length-preserving")
    expanded-row)

  (def result (map mapper lines))
  (assert (= in-dims (dims result)) (printf "expand-row is not dimension-preserving (%q -> %q)" in-dims (dims result)))
  result)

(test (mark-special/expand-row '("")) @(""))
(test (mark-special/expand-row '(".")) @("."))
(test (mark-special/expand-row '("x")) @("x"))
(test (mark-special/expand-row '("..")) @(".."))
(test (mark-special/expand-row '(".x")) @("xx"))
(test (mark-special/expand-row '("x.")) @("xx"))
(test (mark-special/expand-row '("xx")) @("xx"))
(test (mark-special/expand-row '("...")) @("..."))
(test (mark-special/expand-row '("x..")) @("xx."))
(test (mark-special/expand-row '(".x.")) @("xxx"))
(test (mark-special/expand-row '("..x")) @(".xx"))
(test (mark-special/expand-row '("xx.")) @("xxx"))
(test (mark-special/expand-row '(".xx.")) @("xxxx"))
(test (mark-special/expand-row '(".x..xx..x.")) @("xxxxxxxxxx"))

(defn mark-special/expand-all [lines]
  "This function replaces replaces x characters with a 3x3 submatrix."
  (->> lines
       (mark-special/init)
       (mark-special/expand-row)
       (transpose)
       (mark-special/expand-row)
       (transpose)))

########################################
# Functions for producing numerical regions
########################################
(defn numerical-regions/init
  "This function parses a line from a grid into structures that describe the numbers therein."
  [line]
  (defn pair-to-struct [pair]
    (def [num end] pair)
    {:start (- end (length num)) :end end :value (scan-number num)})

  (->> line
       (peg/match '{:main (any (+ (* (<- :d+) ($)) 1))})
       (partition 2)
       (map pair-to-struct)))

(test (numerical-regions/init "...") @[])
(test (numerical-regions/init "467") @[{:start 0 :end 3 :value 467}])
(test (numerical-regions/init "467..114..") @[{:start 0 :end 3 :value 467} {:start 5 :end 8 :value 114}])

(defn numerical-regions/matching [marked-line numerical-regions]
  (defn contains-x [numerical-region]
    (def subline (string/slice marked-line (numerical-region :start) (numerical-region :end)))
    (not (nil? (string/find "x" subline))))
  (filter contains-x numerical-regions))

(test (numerical-regions/matching "..." @[]) @[])
(test (numerical-regions/matching "..." @[{:start 0 :end 3 :value 467}]) @[])
(test (numerical-regions/matching "x.." @[{:start 1 :end 3 :value 46}]) @[])
(test (numerical-regions/matching "x.." @[{:start 0 :end 3 :value 467}]) @[{:start 0 :end 3 :value 467}])
(test (numerical-regions/matching "..xxx....." @[{:start 0 :end 3 :value 467} {:start 5 :end 8 :value 114}]) @[{:start 0 :end 3 :value 467}])

(defn solve [filepath]
  (def data (->> (slurp filepath) (string/trim) (string/split "\n")))
  (def marked-regions (mark-special/expand-all data))
  (def numerical-regions (map numerical-regions/init data))
  (def zipped (zip-with |[$0 $1] marked-regions numerical-regions))

  (defn reducer [acc el]
    (->> (numerical-regions/matching ;el) (map |($ :value)) (sum) (+ acc)))

  (reduce reducer 0 zipped))

(test (solve "integration-part-1") 4361)
(test (solve "input") 543867)
