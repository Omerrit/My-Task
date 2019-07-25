#include <algorithm>
#include <iostream>
#include <vector>
using namespace std;

bool isPalindrome(int t) {
	int reverse = 0;
	if (t < 0) {
		return false;
	}
	else if (t < 10) return true;

	int x = t;
	while (x / 10 != 0)
	{
		reverse *= 10;
		reverse += x % 10;
		x /= 10;
	}
	return reverse == t || reverse == t / 10;
}

int main() {
	while (true) {
		int x;
		cin >> x;
		isPalindrome(x) ? cout << "TRUE" : cout << "False" << endl;
	}
	system("pause");
}