package main

import (
	"errors"
	"fmt"
	"os"
	"time"
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
	WildRank
)

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
	WildSuit
)

type HandValue struct {
	Hand   Hand
	Values []Rank
}

type Card struct {
	Rank Rank
	Suit Suit
}

func cardFromIndex(index int) Card {
	rankIndex := index / 4
	suitIndex := index % 4
	return Card{Rank: Rank(rankIndex), Suit: Suit(suitIndex)}
}

func formatHandPadding(hand Hand) string {
	switch hand {
	case Pair:
		fallthrough
	case Flush:
		return "\t\t\t"
	case ThreeOfAKind:
		return "\t"
	default:
		return "\t\t"
	}
}

func formatHand(hand Hand) string {
	switch hand {
	case HighCard:
		return "High card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two pair"
	case ThreeOfAKind:
		return "Three of a kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full house"
	case FourOfAKind:
		return "Four of a kind"
	case StraightFlush:
		return "Straight flush"
	case RoyalFlush:
		return "Royal flush"
	default:
		return ""
	}
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
	case 'X':
		return WildRank, nil
	}

	return Deuce, errors.New("given character not a valid rank")
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
	case 'x':
		return WildSuit, nil
	}

	return Clubs, errors.New("given character not a valid suit")
}

func inputToCards(input string) ([]Card, error) {
	if len(input)%2 != 0 {
		return nil, errors.New("input length not valid")
	}

	cardAmount := len(input) / 2
	cards := []Card{}

	for i := 0; i < cardAmount; i++ {
		rank, rankError := rankByChar(rune(input[i*2]))

		if rankError != nil {
			return nil, rankError
		}

		suit, suitError := suitByChar(rune(input[i*2+1]))

		if suitError != nil {
			return nil, suitError
		}

		if rank == WildRank && suit != WildSuit ||
			rank != WildRank && suit == WildSuit {
			return nil, errors.New("mixed wild combination not valid")
		}

		if rank != WildRank && suit != WildSuit {
			cards = append(cards, Card{Rank: rank, Suit: suit})
		}
	}

	return cards, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("There must be 1 argument.")
		os.Exit(1)
	}

	input := os.Args[1]

	if len(input) < 15 {
		fmt.Println("Input must have at least length 15.")
		os.Exit(1)
	}

	if (len(input)-15)%5 != 0 {
		fmt.Println("Input length not valid.")
		os.Exit(1)
	}

	boardInput := input[0:10]
	boardCards, boardError := inputToCards(boardInput)

	if boardError != nil {
		fmt.Println("Invalid character(s) found.")
	}

	playersInput := input[10:]
	playerInputs := []string{}
	playerAmount := len(playersInput) / 5

	playersCards := [][]Card{}

	for i := 0; i < playerAmount; i++ {
		playerInput := playersInput[(i*5 + 1):(i*5 + 5)]
		playerCards, playerError := inputToCards(playerInput)

		if playerError != nil {
			fmt.Println("Invalid character(s) found.")
		}

		playerInputs = append(playerInputs, playerInput)
		playersCards = append(playersCards, playerCards)
	}

	wildBoardCards := 5 - len(boardCards)

	wildPlayersCards := 0

	for _, playerCards := range playersCards {
		wildPlayersCards += 2 - len(playerCards)
	}

	var equities [][HandAmount]float64

	start := time.Now()

	if wildBoardCards == 5 && wildPlayersCards == 0 {
		equities = eval_5_2(boardCards, playersCards)
	} else {
		fmt.Println("Format not supported.")
		os.Exit(1)
	}

	elapsed := time.Since(start)

	totalEquity := 0.0

	for _, equity := range equities {
		for _, value := range equity {
			totalEquity += value
		}
	}

	equityPerMillisecond := totalEquity / float64(elapsed.Milliseconds())

	fmt.Printf("Time: %d ms\n", elapsed.Milliseconds())
	fmt.Printf("Speed: %.1f equity/ms\n", equityPerMillisecond)
	fmt.Printf("Total equity: %.1f\n\n", totalEquity)

	for i, equity := range equities {
		playerEquity := 0.0

		for _, value := range equity {
			playerEquity += value
		}

		playerRatioPercent := playerEquity / totalEquity * 100
		fmt.Printf("Player %d - ", i+1)
		fmt.Print(playerInputs[i])
		fmt.Printf("\nTotal:\t\t\t%10.1f %14.8f%%\n", playerEquity, playerRatioPercent)

		for i := 0; i < HandAmount; i++ {
			hand := Hand(i)
			ratioPercent := equity[i] / totalEquity * 100

			fmt.Print(formatHand(hand))
			fmt.Print(":")
			fmt.Print(formatHandPadding(Hand(i)))
			fmt.Printf("%10.1f %14.8f%%\n", equity[i], ratioPercent)
		}

		fmt.Println()
	}
}

func eval_5_2(boardCards []Card, playersCards [][]Card) [][HandAmount]float64 {

	handEquities := [][HandAmount]float64{}

	for range playersCards {
		handEquities = append(handEquities, [HandAmount]float64{})
	}

	deadCards := []Card{}

	for _, playerCards := range playersCards {
		deadCards = append(deadCards, playerCards...)
	}

	aliveCards := []Card{}

	for i := 0; i < CardAmount; i++ {
		card := cardFromIndex(i)
		isDead := false

		for _, deadCard := range deadCards {
			if card.Rank == deadCard.Rank && card.Suit == deadCard.Suit {
				isDead = true
				break
			}
		}

		if !isDead {
			aliveCards = append(aliveCards, card)
		}
	}

	handValues := []HandValue{}

	for range playersCards {
		handValues = append(handValues, HandValue{})
	}

	cardsToEvaluate := [7]Card{}

	for card1 := 0; card1 < len(aliveCards); card1++ {
		for card2 := card1 + 1; card2 < len(aliveCards); card2++ {
			for card3 := card2 + 1; card3 < len(aliveCards); card3++ {
				for card4 := card3 + 1; card4 < len(aliveCards); card4++ {
					for card5 := card4 + 1; card5 < len(aliveCards); card5++ {
						cardsToEvaluate[2] = aliveCards[card1]
						cardsToEvaluate[3] = aliveCards[card2]
						cardsToEvaluate[4] = aliveCards[card3]
						cardsToEvaluate[5] = aliveCards[card4]
						cardsToEvaluate[6] = aliveCards[card5]

						for i, playerCards := range playersCards {
							cardsToEvaluate[0] = playerCards[0]
							cardsToEvaluate[1] = playerCards[1]

							handValues[i] = evalCards(cardsToEvaluate[:])
						}

						winners := []int{0}

						for i := 1; i < len(playersCards); i++ {
							value := handValues[i]
							winnerValue := handValues[winners[0]]

							comparison := int(value.Hand) - int(winnerValue.Hand)

							if comparison == 0 {
								for j := len(value.Values) - 1; j >= 0; j-- {
									comparison = int(value.Values[j]) - int(winnerValue.Values[j])

									if comparison != 0 {
										break
									}
								}
							}

							if comparison > 0 {
								winners = []int{i}
							} else if comparison == 0 {
								winners = append(winners, i)
							}
						}

						winnerEquity := 1.0 / float64(len(winners))

						for i := 0; i < len(winners); i++ {
							winnerIndex := winners[i]
							handIndex := int(handValues[winnerIndex].Hand)
							handEquities[winnerIndex][handIndex] += winnerEquity
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
		values[3] = highestPairRank
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
