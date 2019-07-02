#include<iostream>
using namespace std;
struct digit {
	int value;
	digit(int x) : value(x) {}
};
digit* Sum(digit* a, digit* b) {
	digit* n = new digit(a->value + b->value);
	return n;
}

int main() {
	digit a(5);
	digit b(3);
	auto* l = Sum(a, b);
	std::cout << ;
}