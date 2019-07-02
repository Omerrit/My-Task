#include <vector>
#include <iostream>
#include <numeric>

int main() {
	int n, m;
	std:: cin >> n >> m;
	std::vector<std :: vector<int>> board(n,std :: vector<int>(m));
	for (int i = 0; i < n; ++i) {
		board[i][0] = 1;
	}
	for (int i = 0; i < m; i++) {
		board[0][i] = 1;
	}
	for (int i = 1; i < n; ++i) {
		for (int j = 1; j < m; j++) {
			board[i][j] = board[i - 1][j] + board[i][j - 1];
		}
	}
	std::cout << board[n - 1][m - 1];
	system("pause");
}