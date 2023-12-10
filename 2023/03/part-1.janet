(use judge)

(defn dims [lines]
  (def nrows (length lines))
  (def column-sizes (->> lines (map length) (distinct)))
  (assert (<= (length column-sizes) 1) "the column sizes are non-uniform")
  [nrows (or (column-sizes 0) 0)])

(defn transpose [lines]
  (def in-dims (dims lines))
  (def columns (range 0 (length (lines 0))))
  (defn construct-column [column] (map (fn [line] (line column)) lines))
  (def result (->> columns
                   (map construct-column)
                   (map (fn [row] (string/from-bytes ;row)))))
  (def out-dims (dims result))
  (assert (= [;(reverse in-dims)] out-dims) (printf "not a transposition (dims: %q -> %q)" in-dims out-dims))
  result)

(defn zip-with [f xs ys]
  (def n (min (length xs) (length ys)))
  (reduce (fn [acc el] (array/push acc (f (xs el) (ys el)))) @[] (range n)))

(defn mark-special/init [lines]
  (def pattern ~{:replace-digit (/ (<- :d) ".")
                 :replace-special (/ (<- 1) "x")
                 :main (any (+ :replace-digit (<- ".") :replace-special))})

  (->> lines
       (map |(peg/match pattern $))
       (map string/join)))

(defn mark-special/expand-row [lines]
  (def in-dims (dims lines))
  (def ncols (in-dims 1))

  (def pattern ~{:dot-step (if (+ ".." (* "." -1)) (<- "."))
                 :any-x-or-x-end (if (+ ".x" "xx" (* "x" -1)) (* (constant "x") 1))
                 :x-dot (if "x." (* (constant "xx") 2)) # x.. | x.x | x.$ are all handled.
                 :main (any (+ :dot-step :any-x-or-x-end :x-dot))})

  (defn mapper [row]
    (def expanded-row (string/join (peg/match pattern row)))
    (assert (= (length expanded-row) (length row))
            (printf "the row-expansion is not length-preserving (%d -> %d)\ninput:  %s\noutput: %s" (length row) (length expanded-row) row expanded-row))
    expanded-row)

  (def result (map mapper lines))
  (def out-dims (dims result))
  (assert (= in-dims out-dims) (printf "expand-row is not dimension-preserving (%q -> %q)" in-dims out-dims))
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
  (->> lines
       (mark-special/init)
       (mark-special/expand-row)
       (transpose)
       (mark-special/expand-row)
       (transpose)))

(defn numerical-regions/init [line]
  (defn pair-to-struct [pair]
    (def [num end] pair)
    {:value (scan-number num) :start (- end (length num)) :end end})

  (->> line
       (peg/match ~{:main (any (+ (* (<- :d+) ($)) 1))})
       (partition 2)
       (map pair-to-struct)))

(test (numerical-regions/init "...") @[])
(test (numerical-regions/init "467") @[{:end 3 :start 0 :value 467}])
(test (numerical-regions/init "467..114..") @[{:end 3 :start 0 :value 467} {:end 8 :start 5 :value 114}])

(defn numerical-regions/matching [marked-line numerical-regions]
  (defn contains-x [numerical-region]
    (def xs (string/slice marked-line (numerical-region :start) (numerical-region :end)))
    (not (nil? (peg/match '{:main (* (any (if-not "x" 1)) "x" (any 1))} xs))))
  (filter contains-x numerical-regions))

(test (numerical-regions/matching "..." @[]) @[])
(test (numerical-regions/matching "..." @[{:end 3 :start 0 :value 467}]) @[])
(test (numerical-regions/matching "x.." @[{:end 3 :start 1 :value 46}]) @[])
(test (numerical-regions/matching "x.." @[{:end 3 :start 0 :value 467}]) @[{:end 3 :start 0 :value 467}])
(test (numerical-regions/matching "..xxx....." @[{:end 3 :start 0 :value 467} {:end 8 :start 5 :value 114}]) @[{:end 3 :start 0 :value 467}])

(defn solve [filepath]
  (def data
    (->> (slurp filepath)
         (string/trim)
         (string/split "\n")))

  (def marked-regions (mark-special/expand-all data))
  (def numerical-regions (map numerical-regions/init data))
  (def zipped (zip-with |[$0 $1] marked-regions numerical-regions))

  (defn reducer [acc el]
    (->> (numerical-regions/matching ;el)
         (map |($ :value))
         (sum)
         (+ acc)))

  (reduce reducer 0 zipped))

(test (solve "integration-part-1") 4361)
(test (solve "input") 543867)
