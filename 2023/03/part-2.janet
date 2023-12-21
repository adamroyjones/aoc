(use judge)

######################################
# General functions
######################################
(defn zip-with [f xs ys]
  (def n (min (length xs) (length ys)))
  (reduce (fn [acc el] (array/push acc (f (xs el) (ys el)))) @[] (range n)))

(defn with-index [xs] (zip-with |[$0 $1] xs (range (length xs))))
(defn to-bytes [str] (->> str (string/bytes) (map string/from-bytes)))
(defn point-to-vector [tbl] [(tbl :i) (tbl :j)])

(defn l-infinity [xs ys]
  (assert (= (length xs) (length ys)))
  (->> (length xs)
       (range)
       (map |(math/abs (- (xs $) (ys $))))
       (reduce max 0)))

(test (l-infinity [0 1] [0 1]) 0)
(test (l-infinity [0 0] [0 1]) 1)
(test (l-infinity [1 0] [0 0]) 1)
(test (l-infinity [1 0] [0 1]) 1)
(test (l-infinity [2 0] [0 1]) 2)

######################################
# Matrices
######################################
(defn matrix/new (str)
  (->> str (string/trim) (string/split "\n") (map to-bytes)))

(test (matrix/new "..\nx.") @[@["." "."] @["x" "."]])

(defn matrix/star-locations [matrix]
  (defn reducer [acc row-with-index]
    (def [row row-index] row-with-index)
    (def cells-with-col-indices (with-index row))
    (def star-locations
      (->> cells-with-col-indices
           (filter (fn [cell-with-col-index] (def [cell col-index] cell-with-col-index) (= cell "*")))
           (map (fn [cell-with-col-index] (def [_ col-index] cell-with-col-index) {:i row-index :j col-index}))))
    (array/push acc ;star-locations))

  (reduce reducer @[] (with-index matrix)))

(test (matrix/star-locations @[@["." "."] @["." "."]]) @[])
(test (matrix/star-locations @[@["*" "."] @["." "."]]) @[{:i 0 :j 0}])
(test (matrix/star-locations @[@["." "."] @["." "*"]]) @[{:i 1 :j 1}])
(test (matrix/star-locations @[@["*" "."] @["*" "*"]]) @[{:i 0 :j 0} {:i 1 :j 0} {:i 1 :j 1}])

######################################
# Numerical regions
######################################
(defn numerical-regions/process-line [line]
  (defn pair-to-struct [pair]
    (def [num end] pair)
    @{:start (- end (length num)) :end end :value (scan-number num)})

  (->> line
       (peg/match ~{:main (any (+ (* (<- :d+) ($)) 1))})
       (partition 2)
       (map pair-to-struct)))

(test (numerical-regions/process-line "...") @[])
(test (numerical-regions/process-line "467") @[@{:start 0 :end 3 :value 467}])
(test (numerical-regions/process-line "467..114..") @[@{:start 0 :end 3 :value 467} @{:start 5 :end 8 :value 114}])

(defn numerical-regions/process-lines [lines]
  (defn regions/new [line-with-index]
    (def [line line-index] line-with-index)
    (->> line
         (numerical-regions/process-line)
         (map (fn [region] (put region :i line-index)))))

  (->> lines (with-index) (map regions/new) (flatten)))

(test (numerical-regions/process-lines @["21." "..4"]) @[@{:i 0 :start 0 :end 2 :value 21} @{:i 1 :start 2 :end 3 :value 4}])

(defn numerical-regions/new [str]
  (->> str (string/trim) (string/split "\n") (numerical-regions/process-lines)))

(defn numerical-region/nearest-point [region point]
  (def inclusive-end (dec (region :end)))
  (def inclusive-start (region :start))
  (def row (region :i))
  (def col (point :j))
  (cond (> col inclusive-end) {:i row :j inclusive-end}
        (< col inclusive-start) {:i row :j inclusive-start}
        {:i row :j col}))

(defn numerical-regions/neighbours [numerical-regions star-loc]
  (defn neighbour? [star-loc region]
    (<= (l-infinity (point-to-vector (numerical-region/nearest-point region star-loc)) (point-to-vector star-loc)) 1))

  (filter |(neighbour? star-loc $) numerical-regions))

######################################
# Solving the problem
######################################
(defn solve [filepath]
  (def data (slurp filepath))
  (def star-locations (->> data (matrix/new) (matrix/star-locations)))
  (def regions (numerical-regions/new data))

  (->> star-locations
       (map |(numerical-regions/neighbours regions $))
       (filter |(= (length $) 2))
       (map (fn [regions] (->> regions (map |($ :value)) (product))))
       (sum)))

(test (solve "integration-part-2") 467835)
(test (solve "input") 79613331)
