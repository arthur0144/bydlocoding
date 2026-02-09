package main

import "fmt"

type LinkedList struct {
	head *node
	tail *node
	size int
}

type node struct {
	val  int
	next *node
}

func NewLinkedList() LinkedList {
	return LinkedList{}
}

func (ll *LinkedList) Get(index int) int {
	if index >= ll.size {
		return -1
	}

	if index == ll.size-1 {
		return ll.tail.val
	}

	cur := ll.head
	for i := 0; i < index; i++ {
		cur = cur.next
	}
	return cur.val
}

func (ll *LinkedList) InsertHead(val int) {
	newHead := &node{
		val:  val,
		next: ll.head,
	}

	ll.head = newHead
	ll.size++

	if ll.size == 1 {
		ll.tail = ll.head
	}
}

func (ll *LinkedList) InsertTail(val int) {
	if ll.size == 0 {
		ll.InsertHead(val)
		return
	}

	newTail := &node{val: val}
	ll.tail.next = newTail
	ll.tail = newTail
	ll.size++
}

func (ll *LinkedList) Remove(index int) bool {
	if index >= ll.size {
		return false
	}

	if index == 0 {
		ll.head = ll.head.next
		ll.size--
		if ll.size == 0 {
			ll.tail = nil
		}
		return true
	}

	prev := ll.head
	for i := 0; i < index-1; i++ {
		prev = prev.next
	}
	prev.next = prev.next.next
	if index == ll.size-1 {
		ll.tail = prev
	}
	ll.size--

	return true
}

func (ll *LinkedList) GetValues() []int {
	vals := make([]int, 0, ll.size)
	cur := ll.head
	for range ll.size {
		vals = append(vals, cur.val)
		cur = cur.next
	}
	return vals
}

func main() {
	l := NewLinkedList()
	l.InsertHead(3)
	l.InsertHead(2)
	l.InsertHead(1)
	fmt.Println(l.GetValues())
	l.InsertTail(4)
	fmt.Println(l.GetValues())
	l.Remove(0)
	fmt.Println(l.GetValues())
	l.Remove(2)
	fmt.Println(l.GetValues())
	l.Remove(0)
	fmt.Println(l.GetValues())
	l.Remove(0)
	fmt.Println(l.GetValues())

}
