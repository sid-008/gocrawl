package queue

type Queue struct {
	link []string
}

func (q *Queue) Enqueue(link string) {
	q.link = append(q.link, link)
}

func (q *Queue) Dequeue() string {
	if len(q.link) == 0 {
		return ""
	}
	link := q.link[0]
	q.link = q.link[1:]
	return link
}

func (q *Queue) IsEmpty() bool {
	return len(q.link) == 0
}
