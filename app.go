package main

import (
	"errors"
	"fmt"
	"os"
)

const CardAmount = 52
const RankAmount = 13
const SuitAmount = 4
const HandAmount = 10

type Hand int

const (
	HighCard Hand = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

type Rank int

const (
	Deuce Rank = iota
	Tray
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type HandValue struct {
	Hand   Hand
	Values []Rank
}

type Card struct {
	Rank Rank
	Suit Suit
}

func CardFromIndex(index int) Card {
	rankIndex := index / 4
	suitIndex := index % 4
	return Card{Rank: Rank(rankIndex), Suit: Suit(suitIndex)}
}

func rankByChar(rankChar rune) (Rank, error) {
	switch rankChar {
	case '2':
		return Deuce, nil
	case '3':
		return Tray, nil
	case '4':
		return Four, nil
	case '5':
		return Five, nil
	case '6':
		return Six, nil
	case '7':
		return Seven, nil
	case '8':
		return Eight, nil
	case '9':
		return Nine, nil
	case 'T':
		return Ten, nil
	case 'J':
		return Jack, nil
	case 'Q':
		return Queen, nil
	case 'K':
		return King, nil
	case 'A':
		return Ace, nil
	}

	return Deuce, errors.New("Given character not a valid rank.")
}

func suitByChar(suitChar rune) (Suit, error) {
	switch suitChar {
	case 'c':
		return Clubs, nil
	case 'd':
		return Diamonds, nil
	case 'h':
		return Hearts, nil
	case 's':
		return Spades, nil
	}

	return Clubs, errors.New("Given character not a valid suit.")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("There must be 1 argument.")
		os.Exit(1)
	}

	if os.Args[1] == "XxXxXxXxXx XxXx" {
		equities := eval()
		fmt.Println(equities)
		return
	}

	os.Exit(1)
}

func eval() [HandAmount]int {
	handEquities := [HandAmount]int{}

	for card1 := 0; card1 < CardAmount; card1++ {
		for card2 := card1 + 1; card2 < CardAmount; card2++ {
			for card3 := card2 + 1; card3 < CardAmount; card3++ {
				for card4 := card3 + 1; card4 < CardAmount; card4++ {
					for card5 := card4 + 1; card5 < CardAmount; card5++ {
						for card6 := card5 + 1; card6 < CardAmount; card6++ {
							for card7 := card6 + 1; card7 < CardAmount; card7++ {
								cards := [7]Card{
									CardFromIndex(card1),
									CardFromIndex(card2),
									CardFromIndex(card3),
									CardFromIndex(card4),
									CardFromIndex(card5),
									CardFromIndex(card6),
									CardFromIndex(card7),
								}

								handValue := evalCards(cards[:])
								handEquities[handValue.Hand]++
							}
						}
					}
				}
			}
		}
	}

	return handEquities
}

func evalCards(cards []Card) HandValue {
	// Straight flush

	suitAmounts := [SuitAmount]int{}
	var flushSuit Suit
	flushFound := false

	for _, card := range cards {
		suitAmounts[card.Suit]++

		if suitAmounts[card.Suit] == 5 {
			flushFound = true
			flushSuit = card.Suit
			break
		}
	}

	if flushFound {
		coveredFlushRanks := [RankAmount]bool{}

		for _, card := range cards {
			if card.Suit == flushSuit {
				coveredFlushRanks[card.Rank] = true
			}
		}

		consecutiveAmount := 0

		for i := RankAmount - 1; i >= 0; i-- {
			if coveredFlushRanks[i] {
				consecutiveAmount++
			} else {
				consecutiveAmount = 0
			}

			if consecutiveAmount == 5 {
				if i == int(Ten) {
					return HandValue{
						Hand:   RoyalFlush,
						Values: []Rank{Ace},
					}
				} else {
					return HandValue{
						Hand:   StraightFlush,
						Values: []Rank{Rank(i + 4)},
					}
				}
			} else if consecutiveAmount == 4 && i == int(Deuce) && coveredFlushRanks[Ace] {
				return HandValue{
					Hand:   StraightFlush,
					Values: []Rank{Five},
				}
			}
		}
	}

	rankAmounts := [RankAmount]int{}

	for _, card := range cards {
		rankAmounts[card.Rank]++
	}

	// Four of a kind

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] == 4 {
			for j := RankAmount - 1; j >= 0; j-- {
				if rankAmounts[j] > 0 && j != i {
					return HandValue{
						Hand:   FourOfAKind,
						Values: []Rank{Rank(j), Rank(i)},
					}
				}
			}
		}
	}

	// Full house

	threeOfAKindFound := false
	var threeOfAKindRank Rank

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] == 3 {
			threeOfAKindFound = true
			threeOfAKindRank = Rank(i)
			break
		}
	}

	if threeOfAKindFound {
		for i := RankAmount - 1; i >= 0; i-- {
			if rankAmounts[i] >= 2 && i != int(threeOfAKindRank) {
				return HandValue{
					Hand:   FullHouse,
					Values: []Rank{Rank(i), threeOfAKindRank},
				}
			}
		}
	}

	// Flush

	if flushFound {
		coveredFlushRanks := [RankAmount]bool{}

		for _, card := range cards {
			if card.Suit == flushSuit {
				coveredFlushRanks[card.Rank] = true
			}
		}

		values := [5]Rank{}
		valuesIndex := 4

		for i := RankAmount - 1; i >= 0; i-- {
			if coveredFlushRanks[i] {
				values[valuesIndex] = Rank(i)
				valuesIndex--

				if valuesIndex < 0 {
					break
				}
			}
		}

		return HandValue{
			Hand:   Flush,
			Values: values[:],
		}
	}

	// Straight

	consecutiveAmount := 0

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] > 0 {
			consecutiveAmount++
		} else {
			consecutiveAmount = 0
		}

		if consecutiveAmount == 5 {
			return HandValue{
				Hand:   Straight,
				Values: []Rank{Rank(i + 4)},
			}
		} else if consecutiveAmount == 4 && i == int(Deuce) && rankAmounts[Ace] > 0 {
			return HandValue{
				Hand:   Straight,
				Values: []Rank{Five},
			}
		}
	}

	// Three of a kind

	if threeOfAKindFound {
		values := [3]Rank{}
		values[2] = threeOfAKindRank
		valuesIndex := 1

		for i := RankAmount - 1; i >= 0; i-- {
			if rankAmounts[i] > 0 && i != int(threeOfAKindRank) {
				values[valuesIndex] = Rank(i)
				valuesIndex--

				if valuesIndex < 0 {
					break
				}
			}
		}

		return HandValue{
			Hand:   ThreeOfAKind,
			Values: values[:],
		}
	}

	// Two pair

	pairFound := false
	var highestPairRank Rank

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] == 2 {
			pairFound = true
			highestPairRank = Rank(i)
			break
		}
	}

	secondPairFound := false
	var secondPairRank Rank

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] == 2 && i != int(highestPairRank) {
			secondPairFound = true
			secondPairRank = Rank(i)
			break
		}
	}

	if pairFound && secondPairFound {
		for i := RankAmount - 1; i >= 0; i-- {
			if rankAmounts[i] > 0 && i != int(highestPairRank) && i != int(secondPairRank) {
				return HandValue{
					Hand:   TwoPair,
					Values: []Rank{Rank(i), secondPairRank, highestPairRank},
				}
			}
		}
	}

	if pairFound {
		values := [4]Rank{}
		values[3] = threeOfAKindRank
		valuesIndex := 2

		for i := RankAmount - 1; i >= 0; i-- {
			if rankAmounts[i] > 0 && i != int(highestPairRank) {
				values[valuesIndex] = Rank(i)
				valuesIndex--

				if valuesIndex < 0 {
					break
				}
			}
		}

		return HandValue{
			Hand:   Pair,
			Values: values[:],
		}
	}

	values := [5]Rank{}
	valuesIndex := 4

	for i := RankAmount - 1; i >= 0; i-- {
		if rankAmounts[i] > 0 {
			values[valuesIndex] = Rank(i)
			valuesIndex--

			if valuesIndex < 0 {
				break
			}
		}
	}

	return HandValue{
		Hand:   HighCard,
		Values: values[:],
	}
}
