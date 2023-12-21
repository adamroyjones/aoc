(use judge)

(def constraints {:red 12 :green 13 :blue 14})

(defn round-possible? [str]
  (def pattern ~{:main (* :s* (/ (<- :d*) ,scan-number) :s* (/ (<- :a*) ,keyword))})
  (->> str
       (string/split ",")
       (map (fn [round] (peg/match pattern round)))
       (map (fn [pair] (def [count colour] pair) (<= count (constraints colour))))
       (every?)))

(defn game-possible? [str]
  (->> str (string/split ";") (map round-possible?) (every?)))

(test (game-possible? "3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green") true)
(test (game-possible? "1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue") true)
(test (game-possible? "8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red") false)
(test (game-possible? "1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red") false)
(test (game-possible? "6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green") true)

(defn parse-line [line]
  (def pattern ~{:main (* "Game " (/ (<- :d*) ,scan-number) ":" (/ (<- (any 1)) ,game-possible?))})
  (peg/match pattern line))

(defn solve [filename]
  (defn reducer [acc el]
    (def [id possible] (parse-line el))
    (if possible (+ acc id) acc))

  (with [fd (file/open filename)]
        (reduce reducer 0 (file/lines fd))))

(test (solve "integration-part-1") 8)
(test (solve "input") 2265)
