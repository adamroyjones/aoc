(use judge)

######################################
# General functions
######################################
(defn zip-with [f xs ys]
  (def n (min (length xs) (length ys)))
  (reduce (fn [acc el] (array/push acc (f (xs el) (ys el)))) @[] (range n)))

(defn to-bytes [str] (->> str (string/bytes) (map string/from-bytes)))

(defn point-to-vector [tbl] [(tbl :i) (tbl :j)])

(defn l-infinity [xs ys]
  (assert (= (length xs) (length ys)))
  (->> (range 0 (length xs))
       (map |(math/abs (- (xs $) (ys $))))
       (reduce (fn [acc el] (max acc el)) 0)))

(test (l-infinity [0 1] [0 1]) 0)
(test (l-infinity [0 0] [0 1]) 1)
(test (l-infinity [1 0] [0 0]) 1)
(test (l-infinity [1 0] [0 1]) 1)
(test (l-infinity [2 0] [0 1]) 2)

######################################
# Matrix functions
######################################
(defn matrix (str)
  (->> str
       (string/trim)
       (string/split "\n")
       (map to-bytes)))

(test (matrix "..\nx.") @[@["." "."] @["x" "."]])

(defn matrix-to-star-locations [matrix]
  (def matrix-with-indices (zip-with |[$0 $1] matrix (range 0 (length matrix))))

  (defn reducer [locations row-with-index]
    (def [row row-index] row-with-index)
    (def cells-with-col-indices (zip-with |[$0 $1] row (range 0 (length row))))
    (def matching-cell-indices
      (->> cells-with-col-indices
           (filter (fn [cell-with-index] (def [cell col-index] cell-with-index) (= cell "*")))
           (map (fn [cell-with-index] (def [_ col-index] cell-with-index) {:i row-index :j col-index}))))
    (array/push locations ;matching-cell-indices))

  (reduce reducer @[] matrix-with-indices))

(test (matrix-to-star-locations @[@["*" "."] @["*" "*"]]) @[{:i 0 :j 0} {:i 1 :j 0} {:i 1 :j 1}])

######################################
# Numerical region functions
######################################
(defn numerical-regions/process-row [row]
  (defn pair-to-struct [pair]
    (def [num end] pair)
    @{:value (scan-number num) :start (- end (length num)) :end end})

  (->> row
       (peg/match ~{:main (any (+ (* (<- :d+) ($)) 1))})
       (partition 2)
       (map pair-to-struct)))

(test (numerical-regions/process-row "...") @[])
(test (numerical-regions/process-row "467") @[@{:end 3 :start 0 :value 467}])
(test (numerical-regions/process-row "467..114..") @[@{:end 3 :start 0 :value 467} @{:end 8 :start 5 :value 114}])

(defn numerical-regions/process-rows [rows]
  (def rows-with-row-indices (zip-with |[$0 $1] rows (range 0 (length rows))))

  (defn row-with-row-index-to-regions [row-with-row-index]
    (def [row row-index] row-with-row-index)
    (->> (numerical-regions/process-row row)
         (map (fn [region] (put region :i row-index)))))

  (->> rows-with-row-indices
       (map row-with-row-index-to-regions)
       (flatten)))

(test (numerical-regions/process-rows @["21." "..4"]) @[@{:end 2 :i 0 :start 0 :value 21} @{:end 3 :i 1 :start 2 :value 4}])

(defn numerical-regions [str]
  (->> str
       (string/trim)
       (string/split "\n")
       (numerical-regions/process-rows)))

######################################
# Solving the problem
######################################
(defn numerical-region/nearest-point [region point]
  (def inclusive-end (dec (region :end)))
  (def inclusive-start (region :start))
  (def row (region :i))
  (def col (point :j))
  (cond (> col inclusive-end) {:i row :j inclusive-end}
        (< col inclusive-start) {:i row :j inclusive-start}
        {:i row :j col}))

(defn star-location-to-matching-numerical-regions [star-location numerical-regions]
  (def regions-with-nearest-points
    (map (fn [region] {:region region :nearest-point (numerical-region/nearest-point region star-location)})
         numerical-regions))

  (defn nearby-regions-filter [region-with-nearest-point]
    (def {:region region :nearest-point nearest-point} region-with-nearest-point)
    (<= (l-infinity (point-to-vector nearest-point) (point-to-vector star-location)) 1))

  (def nearby-regions (filter nearby-regions-filter regions-with-nearest-points))
  (assert (<= (length nearby-regions) 2))
  (map (fn [region-with-nearest-point] (region-with-nearest-point :region)) nearby-regions))

(defn solve [filepath]
  (def data (slurp filepath))
  (def star-locations (->> data (matrix) (matrix-to-star-locations)))
  (def regions (numerical-regions data))

  (->> star-locations
       (map (fn [star-location] (star-location-to-matching-numerical-regions star-location regions)))
       (filter (fn [regions] (= (length regions) 2)))
       (map (fn [regions] (->> regions (map |($ :value)) (product))))
       (sum)))

(test (solve "integration-part-2") 467835)
(test (solve "input") 79613331)
