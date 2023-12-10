(use judge)

(def constraints {:red 12 :green 13 :blue 14})

(defn min-constraint [rounds]
  (def min-hash @{:red 0 :green 0 :blue 0})
  (each round rounds
    (each pair round
      (def [count colour] pair)
      (put min-hash colour (max (min-hash colour) count))))
  min-hash)

(defn parse-round [str]
  (def pattern ~{:main (* :s* (/ (<- :d*) ,scan-number) :s* (/ (<- :a*) ,keyword))})
  (->> str
       (string/split ",")
       (map (fn [round] (peg/match pattern round)))))

(defn power [str]
  (->> str
       (string/split ":")
       (last)
       (string/split ";")
       (map parse-round)
       (min-constraint)
       (values)
       (product)))

(test (power "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green") 48)
(test (power "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue") 12)
(test (power "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red") 1560)
(test (power "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red") 630)
(test (power "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green") 36)

(defn solve [filename]
  (with [fd (file/open filename)]
        (->> (file/lines fd) (map power) (sum))))

(test (solve "integration-part-2") 2286)
(test (solve "input") 64097)
