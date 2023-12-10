(use judge)

(defn array-to-table [xs] (reduce (fn [tbl x] (put tbl x true)) @{} xs))

(test (array-to-table []) @{})
(test (array-to-table [1]) @{1 true})
(test (array-to-table [1 2]) @{1 true 2 true})
(test (array-to-table [1 1]) @{1 true})

(defn parse-line [line]
  (def pattern ~{:game (* "Card" :s* :d* ":")
                 :numbers (any (+ :s+ (/ (<- :d+) ,scan-number)))
                 :main (* :game (group :numbers) "|" (group :numbers))})

  (peg/match pattern line))

(test (parse-line "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53") @[@[41 48 83 86 17] @[83 86 6 31 17 9 48 53]])

(defn score [parsed-line]
  (assert (= (length parsed-line) 2))
  (def winning (array-to-table (parsed-line 0)))
  (def ticket-numbers (parsed-line 1))
  (def count (reduce (fn [acc ticket-number] (if (nil? (get winning ticket-number)) acc (inc acc))) 0 ticket-numbers))
  (if (zero? count) 0 (math/exp2 (dec count))))

(defn solve [filename]
  (->> filename
       (slurp)
       (string/trim)
       (string/split "\n")
       (map parse-line)
       (map score)
       (sum)))

(test (solve "integration-part-1") 13)
(test (solve "input") 23441)
