#include<string>
#include <iostream>
using namespace std;
class Solution {
public:
    int lengthOfLongestSubstring(string s) {
		vector m(128, -1);
        int res = 0, left = -1;
            for (int i = 0; i < s.size(); ++i) {
				left = max(left, m[s[i]]);
				m[s[i]] = i;
				res = max(res, i - left);
			}
		return res;
    }
};

int main() {
	Solution d;
	std::cout << d.lengthOfLongestSubstring("sddkjodwdld");
	system("pause");
}