#include <string>
#include <vector>
#include <iostream>

using namespace std;
void recursion(vector<string>& solutions, string solution, int n) {
	if (solution.size() == 2 * n)
		solutions.push_back(solution);
	else {
		int left = 0;
		int right = 0;
		for (const auto& s : solution) {
			(s == '(') ? ++left : ++right;
		}
		if (left > right) {
			if (left < n)
				recursion(solutions, solution + '(', n);
			recursion(solutions, solution + ')', n);
		} else 
			recursion(solutions, solution + '(', n); 
	}
}

int main() {
	vector<string> solutions;
	int n;
	cin >> n;
	recursion(solutions, "", n);
	for (const auto& s : solutions) {
		for (const auto& p : s) {
			cout << p;
		}
		cout << endl;
	}
	system("pause");
}
