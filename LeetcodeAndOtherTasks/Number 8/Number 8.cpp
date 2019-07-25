#include <algorithm>
#include <iostream>
#include <string>
#include <sstream>
using namespace std;

int myAtoi(string str) {
	stringstream s("");
	if (count_if(str.begin(), str.end(), [](char v) {return v != ' '; }) == 0) {
		return 0;
	}
	s << str;
	int k;
	s >> k;
	return k;
}
int main() {
	string s("  ");
	cout << myAtoi(s);
	system("pause");
}