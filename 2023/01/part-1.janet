(use judge)

(defn parse-line [line]
  (def pattern ~{:main (any (+ (<- :d) 1))})
  (def digits (peg/match pattern line))
  (scan-number (string (first digits) (last digits))))

(test (parse-line "1abc2") 12)
(test (parse-line "pqr3stu8vwx") 38)
(test (parse-line "a1b2c3d4e5f") 15)
(test (parse-line "treb7uchet") 77)

(defn solve [filename]
  (with [fd (file/open filename)]
        (reduce (fn [acc el] (+ acc (parse-line el))) 0 (file/lines fd))))

(test (solve "integration-part-1") 142)
(test (solve "input") 54597)
