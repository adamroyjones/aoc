(use judge)

(defn set/new [xs] (reduce (fn [tbl x] (put tbl x true)) @{} xs))

(test (set/new []) @{})
(test (set/new [1]) @{1 true})
(test (set/new [1 2]) @{1 true 2 true})
(test (set/new [1 1]) @{1 true})

(defn parse-line [line]
  (def pattern ~{:game (* "Card" :s* :d* ":")
                 :numbers (any (+ :s+ (/ (<- :d+) ,scan-number)))
                 :main (* :game (group :numbers) "|" (group :numbers))})

  (peg/match pattern line))

(test (parse-line "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53") @[@[41 48 83 86 17] @[83 86 6 31 17 9 48 53]])

(defn score [parsed-line]
  (assert (= (length parsed-line) 2))
  (def winning-numbers (set/new (first parsed-line)))
  (def ticket-numbers (last parsed-line))
  (def matches (count |(not (nil? (get winning-numbers $))) ticket-numbers))
  (if (zero? matches) 0 (math/exp2 (dec matches))))

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
