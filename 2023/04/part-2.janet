(use judge)

######################################
# General functions
######################################
(defn array-to-table [xs] (reduce (fn [tbl x] (put tbl x true)) @{} xs))

(test (array-to-table []) @{})
(test (array-to-table [1]) @{1 true})
(test (array-to-table [1 2]) @{1 true 2 true})
(test (array-to-table [1 1]) @{1 true})

(defn repeat [val cnt] (->> (range cnt) (map (fn [_] val))))

######################################
# Parsing and counting
######################################
(defn parse-line [line]
  (def pattern ~{:game (* "Card" :s* (/ (<- :d*) ,(comp dec scan-number)) ":") # Decrementing to avoid off-by-one pain.
                 :numbers (any (+ :s+ (/ (<- :d+) ,scan-number)))
                 :main (* :game (group :numbers) "|" (group :numbers))})

  (peg/match pattern line))

(test (parse-line "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53") @[0 @[41 48 83 86 17] @[83 86 6 31 17 9 48 53]])

(defn count [parsed-line]
  (assert (= (length parsed-line) 3))
  (def winning (array-to-table (parsed-line 1)))
  (def ticket-numbers (parsed-line 2))
  (reduce (fn [acc ticket-number] (if (nil? (get winning ticket-number)) acc (inc acc))) 0 ticket-numbers))

(test (count (parse-line "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53")) 4)
(test (count (parse-line "Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19")) 2)
(test (count (parse-line "Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1")) 2)
(test (count (parse-line "Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83")) 1)
(test (count (parse-line "Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36")) 0)
(test (count (parse-line "Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11")) 0)

######################################
# Solving the problem
######################################
(defn state-reducer [cur-state parsed-line]
  (assert (= (length parsed-line) 3))

  (def position (parsed-line 0))
  (def card-count (cur-state position))
  (def max-step (count parsed-line))
  (def steps (range 1 (inc max-step)))

  (each step steps
    (def stepped-position (+ position step))
    (when (< stepped-position (length cur-state))
      (def prev-count (cur-state stepped-position))
      (put cur-state stepped-position (+ prev-count card-count))))

  cur-state)

(defn solve [filename]
  (def parsed-lines (->> filename
                         (slurp)
                         (string/trim)
                         (string/split "\n")
                         (map parse-line)))

  (def initial-state (repeat 1 (length parsed-lines)))
  (def final-state (reduce state-reducer initial-state parsed-lines))
  (sum final-state))

(test (solve "integration-part-1") 30)
(test (solve "input") 5923918)
