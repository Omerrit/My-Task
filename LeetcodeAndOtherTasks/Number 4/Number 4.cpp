#include<algorithm>
#include <iostream>
#include <vector>
using namespace std;
class Solution {
public:
    double findMedianSortedArrays(vector<int>& nums1, vector<int>& nums2) {
    vector<int> nums3;
	merge(nums1.begin(), nums1.end(), nums2.begin(), nums2.end(), back_inserter(nums3));
	int middle = nums3.size() / 2 ;
	if (nums3.size() % 2 == 0) {
		return nums3[middle - 1] + static_cast<double>(nums3[middle] - nums3[middle - 1]) / 2 ;
	} else {
		return nums3[middle];
	}
    }
};

int main() {
	Solution d;
	std::cout << d.lengthOfLongestSubstring("sddkjodwdld");
	system("pause");
}