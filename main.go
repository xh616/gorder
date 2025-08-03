package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func MidOrderStack(root *TreeNode) {
	var stack []*TreeNode
	for len(stack) > 0 || root != nil {
		// 遍历到最左子节点
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		// 此时 root 为 nil，栈顶是最后一个左子节点
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1] // 出栈
		fmt.Println(root.Val)
		// 转向右子树
		root = root.Right
	}
}

func PreOrderStack(root *TreeNode) {
	if root == nil {
		return
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		// 栈顶元素出栈
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Println(root.Val)
		// 先入栈右子树，再入栈左子树，这样出栈时，先左后右
		if root.Right != nil {
			stack = append(stack, root.Right)
		}
		if root.Left != nil {
			stack = append(stack, root.Left)
		}
	}
}

func PostOrderStack(root *TreeNode) {
	if root == nil {
		return
	}
	var stack []*TreeNode
	var prev *TreeNode // 记录上一次访问的节点
	cur := root
	for len(stack) > 0 || cur != nil {
		// 遍历到最左子节点
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		// 栈顶是最后一个左子节点，先不出栈，而是转向右子树
		top := stack[len(stack)-1]
		// 关键判断：如果右子节点存在且未被访问过，优先处理右子树
		if top.Right != nil && prev != top.Right {
			cur = top.Right // 转向右子树，继续循环
		} else {
			// 右子树已处理或不存在，处理当前节点
			fmt.Println(top.Val)
			prev = top                   // 记录已访问的节点
			stack = stack[:len(stack)-1] // 出栈
			// 此时cur为上一次cur的值，栈顶是当前节点的父节点
		}
	}
}

func main() {

}
