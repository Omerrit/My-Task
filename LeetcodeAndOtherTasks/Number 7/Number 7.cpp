#include<string>
#include <iostream>
using namespace std;
class Solution {
public:
	int reverse(int x) {
		std::string t(std::to_string(x));
		auto begin = t.begin();
		if (t[0] == '-') {
			begin++;
		}
		std::reverse(begin, t.end());
		char *end;
		int64_t n = strtol(t.c_str(), &end, 10);
		if (n >= INT_MAX || n <= INT_MIN) {
			return 0;
		}
		return static_cast<int>(n);
	}
};

int main() {
	Solution d;
	std::cout << d.reverse(12);
	system("pause");
}