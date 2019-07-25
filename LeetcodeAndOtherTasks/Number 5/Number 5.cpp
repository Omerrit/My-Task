#include <algorithm>
#include <iostream>
#include <string>
using namespace std;

string longestPalindrome(string s) {
	{

		if (s.length() < 2) return s;
		int PolBegin = 0;
		int maxLen = 1;

		int i = 0;
		while (i < s.length()) {
			int begin = i;
			int end = i;

			while (end + 1 < s.length() && s[end] == s[end + 1]) { end++; }

			i = end + 1;

			int left = begin;
			int right = end;

			while (left - 1 >= 0 && right + 1 < s.length() && s[left - 1] == s[right + 1]) {
				left--;
				right++;
			}

			int length = right - left + 1;
			if (length > maxLen) {
				PolBegin = left;
				maxLen = length;
			}
		}

		return s.substr(PolBegin, maxLen);
	}

}

int main() {
	string s("cacd");
	cout << longestPalindrome(s);
	system("pause");
}