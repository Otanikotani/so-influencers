package main

import (
	"strconv"
	"time"
)

func questionCsv(question *Question) []string {
	return []string{
		questionVerticeId(question),
		"Question",
		question.Title,
		strconv.Itoa(question.ViewCount),
		strconv.Itoa(question.AnswerCount),
		strconv.Itoa(question.Score),
		strconv.FormatBool(question.IsAnswered),
		neptuneDate(question.CreationDate),
	}
}

func answerCsv(answer *Answer) []string {
	return []string{
		answerVerticeId(answer),
		"Answer",
		answer.Title,
		strconv.FormatBool(answer.IsAccepted),
		strconv.Itoa(answer.Score),
		neptuneDate(answer.CreationDate),
	}
}

func shallowUserCsv(shallowUser *ShallowUser) []string {
	return []string{
		shallowUserVerticeId(shallowUser),
		"Person",
		shallowUser.DisplayName,
		strconv.Itoa(shallowUser.Reputation),
	}
}

func shallowUserVerticeId(shallowUser *ShallowUser) string {
	return "u" + strconv.Itoa(shallowUser.UserId)
}

func answerVerticeId(answer *Answer) string {
	return "a" + strconv.Itoa(answer.AnswerId)
}

func questionVerticeId(question *Question) string {
	return "q" + strconv.Itoa(question.QuestionId)
}

func neptuneDate(timestamp int64) string {
	unixDate := time.Unix(timestamp, 0)
	return unixDate.Format("2006-01-02T15:04:05") //YYYY-MM-DDTHH:mm:SS
}

func questionVertices(questions *[]Question) [][]string {
	var questionVertices [][]string
	questionVerticesHeader := []string{"~id", "~label", "title:String", "viewCount:Int", "answerCount:Int", "score:Int", "isAnswered:Bool", "creationDate:Date"}
	questionVertices = append(questionVertices, questionVerticesHeader)

	for _, question := range *questions {
		questionCsvRow := questionCsv(&question)
		questionVertices = append(questionVertices, questionCsvRow)
	}
	return questionVertices
}

func answerVertices(questions *[]Question) [][]string {
	var answerVertices [][]string
	answerVerticesHeader := []string{"~id", "~label", "title:String", "accepted:Bool", "score:Int", "creationDate:Date"}
	answerVertices = append(answerVertices, answerVerticesHeader)

	for _, question := range *questions {
		for _, answer := range question.Answers {
			answerCsvRow := answerCsv(&answer)
			answerVertices = append(answerVertices, answerCsvRow)
		}
	}
	return answerVertices
}

func peopleVertices(questions *[]Question) [][]string {
	var peopleVertices [][]string
	peopleVerticesHeader := []string{"~id", "~label", "title:DisplayName", "reputation:Int"}
	peopleVertices = append(peopleVertices, peopleVerticesHeader)

	peopleByIds := make(map[int]ShallowUser)

	for _, question := range *questions {
		peopleByIds[question.Owner.UserId] = question.Owner
		for _, answer := range question.Answers {
			peopleByIds[answer.Owner.UserId] = answer.Owner
		}
	}

	for _, user := range peopleByIds {
		userCsvRow := shallowUserCsv(&user)
		peopleVertices = append(peopleVertices, userCsvRow)
	}

	return peopleVertices
}

func edges(questions *[]Question) [][]string {
	var edges [][]string
	edgeHeader := []string{"~id", "~from", "~to", "~label"}
	edges = append(edges, edgeHeader)
	edgeCount := 0

	for _, question := range *questions {
		edgeCount++
		edgeId := "e" + strconv.Itoa(edgeCount)
		askedEdge := []string{edgeId, shallowUserVerticeId(&question.Owner), questionVerticeId(&question), "Asked"}
		edges = append(edges, askedEdge)

		for _, answer := range question.Answers {
			edgeCount++
			edgeId := "e" + strconv.Itoa(edgeCount)
			answerEdge := []string{edgeId, shallowUserVerticeId(&answer.Owner), answerVerticeId(&answer), "Answered"}
			edges = append(edges, answerEdge)
		}
	}

	return edges
}

func toVerticesAndEdges(questions *[]Question) ([][]string, [][]string, [][]string, [][]string) {
	return questionVertices(questions), answerVertices(questions), peopleVertices(questions), edges(questions)
}
