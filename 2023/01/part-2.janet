(use judge)

(defn parse-line [line]
  (def number-words '("one" "two" "three" "four" "five" "six" "seven" "eight" "nine"))
  (def fst-pattern ~{:main (any (+ (<- (+ :d ,;number-words)) 1))})
  (def fst (first (peg/match fst-pattern line)))
  (def reversed-number-words (map string/reverse number-words))
  (def snd-pattern ~{:main (any (+ (<- (+ :d ,;reversed-number-words)) 1))})
  (def snd (->> line (string/reverse) (peg/match snd-pattern) (first) (string/reverse)))
  (def kvs (zipcoll number-words (range 1 (+ (length number-words) 1))))
  (defn to-num [digit] (or (kvs digit) (scan-number digit)))
  (+ (* 10 (to-num fst)) (to-num snd)))

(test (parse-line "two1nine") 29)
(test (parse-line "eightwothree") 83)
(test (parse-line "abcone2threexyz") 13)
(test (parse-line "xtwone3four") 24)
(test (parse-line "4nineeightseven2") 42)
(test (parse-line "zoneight234") 14)
(test (parse-line "7pqrstsixteen") 76)
# This trips up a more straightforward attempt to parse lines.
(test (parse-line "fourxzhgjfrrbmkcheightfive7seven8oneightb") 48)

(defn solve [filename]
  (with [fd (file/open filename)]
        (reduce (fn [acc el] (+ acc (parse-line el))) 0 (file/lines fd))))

(test (solve "integration-part-2") 281)
(test (solve "input") 54504)
