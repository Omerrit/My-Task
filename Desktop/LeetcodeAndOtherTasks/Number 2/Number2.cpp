#include<string>
#include<utility>
using namespace std;
 struct ListNode {
     int val;
     ListNode *next;
     ListNode(int x) : val(x), next(NULL) {}
 };
 
 class Solution {
 private:
	 pair<int, int> Sum(const int& a, const int& b) {
		 int value = a + b;
		 if (value >= 10) {
			 return { value % 10, 1 };
		 }
		 else
			 return { value, 0 };
	 }
 public:
	 ListNode* addTwoNumbers(ListNode* l1, ListNode* l2) {
		 pair<int, int> sum = Sum(l1->val, l2->val);
		 ListNode* l3 = new ListNode(sum.first);
		 ListNode* L = new ListNode(0);
		 l3->next = L;
		 int next_plus_one = sum.second;
		 while (l1->next != NULL && l2->next != NULL) {
			 if (l1->next != NULL) {
				 l1 = l1->next;
			 }
			 if (l2->next != NULL) {
				 l2 = l2->next;
			 }
			 sum = Sum(l1->val, l2->val + next_plus_one);
			 L->val = sum.first;
			 next_plus_one = sum.second;
			 L->next = new ListNode(0);
			 L = L->next();
		 }
		 return l3;
	 }
 };

 int main() {
 }